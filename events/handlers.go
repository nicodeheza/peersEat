package events

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
)

type Handlers struct {
	peerRepo   repositories.PeerRepositoryI
	validation validations.ValidateI
	geo        geo.GeoServiceI
}

type HandlersI interface {
	HandleAddPeer(event types.Event)
}

func NewEventHandlers(
	peerRepo repositories.PeerRepositoryI,
	validation validations.ValidateI,
	geo geo.GeoServiceI,
) *Handlers {
	return &Handlers{peerRepo, validation, geo}
}

func (h *Handlers) createSendMap(urls []string, sendMap map[string][]string) {
	if len(urls) == 0 {
		return
	}
	if len(urls) == 1 {
		sendMap[urls[0]] = nil
		return
	}
	cutIndex := int(len(urls) / 2)
	list1 := urls[0:cutIndex]
	list2 := urls[cutIndex:]

	sendMap[list1[0]] = list1[1:]
	sendMap[list2[0]] = list2[1:]
}

func (h *Handlers) sendEvent(peerUrl string, event types.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	url := peerUrl + "/peer/event"

	body, err := json.Marshal(event)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println(err.Error())
		return
	}

	if resp.StatusCode != 200 && len(event.SendTo) > 0 {
		// add foul
		log.Printf("failed to send new peer to %v, retraining with %v\n", peerUrl, event.SendTo[0])
		newUrl := event.SendTo[0]
		event.SendTo = event.SendTo[1:]
		wg.Add(1)
		h.sendEvent(newUrl, event, wg)
		return
	}

	if resp.StatusCode != 200 && len(event.SendTo) == 0 {
		log.Println("Send to is empty")
		return
	}

}

func (h *Handlers) PropagateEvent(event types.Event) {
	if len(event.SendTo) == 0 {
		return
	}
	sendMap := make(map[string][]string)
	h.createSendMap(event.SendTo, sendMap)

	var wg sync.WaitGroup
	for url, sendTo := range sendMap {
		wg.Add(1)
		eventToSend := event
		eventToSend.SendTo = sendTo
		go h.sendEvent(url, eventToSend, &wg)
	}

	go func() {
		wg.Wait()
	}()
}

func (h *Handlers) HandleAddPeer(event types.Event) {
	defer h.PropagateEvent(event)

	if event.Name != "addPeer" {
		return
	}

	selfPeer, err := h.peerRepo.GetSelf()
	if err != nil {
		log.Println(err.Error())
		return
	}

	newPeer, ok := event.Payload.(models.Peer)
	errors := h.validation.ValidatePeer(newPeer)
	if !ok || errors != nil {
		log.Println("payload don't contains a peer")
		return
	}

	id, err := h.peerRepo.Insert(newPeer)
	if err != nil {
		log.Println(err.Error())
		return
	}

	updatedFields := []string{}
	if h.geo.AreInfluenceAreasOverlaying(selfPeer, newPeer) {
		selfPeer.InAreaPeers = append(selfPeer.InAreaPeers, id)
		updatedFields = append(updatedFields, "InAreaPeers")
	}

	if h.geo.IsInDeliveryArea(selfPeer, newPeer) {
		selfPeer.InDeliveryAreaPeers = append(selfPeer.InDeliveryAreaPeers, id)
		updatedFields = append(updatedFields, "InDeliveryAreaPeers")
	}

	if len(updatedFields) > 0 {
		err := h.peerRepo.Update(selfPeer, updatedFields)

		if err != nil {
			log.Println(err.Error())
			return
		}
	}

}
