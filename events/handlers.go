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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handlers struct {
	peerRepo   repositories.PeerRepositoryI
	validation validations.ValidateI
	geo        geo.GeoServiceI
}

type HandlersI interface {
	PropagateEvent(event types.Event)
	HandleAddPeer(event types.Event)
	PeerUpdatedDeliveryArea(event types.Event)
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

	if event.Name != ADD_NEW_PEER {
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
		updatedFields = append(updatedFields, "in_area_peers")
	}

	if h.geo.IsInDeliveryArea(selfPeer, newPeer) {
		selfPeer.InDeliveryAreaPeers = append(selfPeer.InDeliveryAreaPeers, id)
		updatedFields = append(updatedFields, "in_area_delivery_peers")
	}

	if len(updatedFields) > 0 {
		err := h.peerRepo.Update(selfPeer, updatedFields)

		if err != nil {
			log.Println(err.Error())
			return
		}
	}

}

func (h *Handlers) PeerUpdatedDeliveryArea(event types.Event) {
	defer h.PropagateEvent(event)

	if event.Name != DELIVERY_AREA_UPDATED {
		return
	}

	selfPeer, err := h.peerRepo.GetSelf()
	if err != nil {
		log.Println(err.Error())
		return
	}

	sendPeer, ok := event.Payload.(models.Peer)
	errors := h.validation.ValidatePeer(sendPeer)
	if !ok || errors != nil {
		log.Println("payload don't contains a peer")
		return
	}

	peer, err := h.peerRepo.FindByUrlAndUpdate(sendPeer.Url, map[string]interface{}{
		"delivery_radius": sendPeer.DeliveryRadius,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	isInDeliveryAres := h.geo.IsInDeliveryArea(selfPeer, peer)

	var isInDeliveryAreaSlice bool
	for _, id := range selfPeer.InDeliveryAreaPeers {
		if id == peer.Id {
			isInDeliveryAreaSlice = true
			break
		}
	}

	var selfChanged bool

	if isInDeliveryAres && !isInDeliveryAreaSlice {
		selfPeer.InDeliveryAreaPeers = append(selfPeer.InDeliveryAreaPeers, peer.Id)
		selfChanged = true
	}
	if !isInDeliveryAres && isInDeliveryAreaSlice {
		newDeliverySlice := []primitive.ObjectID{}
		for _, id := range selfPeer.InDeliveryAreaPeers {
			if id != peer.Id {
				newDeliverySlice = append(newDeliverySlice, id)
			}
		}

		selfPeer.InDeliveryAreaPeers = newDeliverySlice
		selfChanged = true
	}

	if selfChanged {
		err = h.peerRepo.Update(selfPeer, []string{"in_area_delivery_peers"})
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

}
