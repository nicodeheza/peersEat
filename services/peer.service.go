package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/types"
)

type PeerServiceI interface {
	InitPeer()
	AddNewPeer(newPeer models.Peer) error
	GetSendMap(urls []string, sendMap map[string][]string)
	GetNewSendMap(excludes []string, sendMap map[string][]string) error
	SendNewPeer(body types.PeerPresentationBody, peerUrl string, ch chan<- error, wg *sync.WaitGroup)
	AllPeersToSend(excludeUrls []string) ([]models.Peer, error)
}

type PeerService struct {
	repo repositories.PeerRepositoryI
	geo  geo.GeoServiceI
}

func NewPeerService(repository repositories.PeerRepositoryI, geo geo.GeoServiceI) *PeerService {
	return &PeerService{repository, geo}
}

func (p *PeerService) InitPeer() {
	centerSrt := os.Getenv("CENTER")
	centerSlice := strings.Split(centerSrt, ",")
	long, _ := strconv.ParseFloat(centerSlice[0], 64)
	lat, _ := strconv.ParseFloat(centerSlice[1], 64)

	selfPeer := models.Peer{
		Url:     os.Getenv("HOST"),
		Center:  models.Center{Long: long, Lat: lat},
		City:    os.Getenv("CITY"),
		Country: os.Getenv("COUNTRY"),
	}

	p.repo.Insert(selfPeer)

	// Todo: create new peers endpoint and add here call
}

func (p *PeerService) AddNewPeer(newPeer models.Peer) error {
	id, err := p.repo.Insert(newPeer)
	if err != nil {
		return err
	}

	selfPeer, err := p.repo.GetSelf()
	if err != nil {
		return err
	}

	if newPeer.Country != selfPeer.Country || newPeer.City != selfPeer.City {
		return nil
	}

	// do it concurrent?
	updatedFields := []string{}
	if p.geo.IsInInfluenceArea(selfPeer, newPeer) {
		selfPeer.InAreaPeers = append(selfPeer.InAreaPeers, id)
		updatedFields = append(updatedFields, "InAreaPeers")
	}

	if p.geo.IsInDeliveryArea(selfPeer, newPeer) {
		selfPeer.InDeliveryAreaPeers = append(selfPeer.InDeliveryAreaPeers, id)
		updatedFields = append(updatedFields, "InDeliveryAreaPeers")
	}

	if len(updatedFields) > 0 {
		err := p.repo.Update(selfPeer, updatedFields)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PeerService) GetSendMap(urls []string, sendMap map[string][]string) {
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

func (p *PeerService) GetNewSendMap(excludes []string, sendMap map[string][]string) error {
	allUrls, err := p.repo.GetAllUrls(excludes)

	if err != nil {
		return err
	}

	p.GetSendMap(allUrls, sendMap)
	return nil
}

func (p *PeerService) SendNewPeer(body types.PeerPresentationBody, peerUrl string, ch chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	url := peerUrl + "/peer/present"

	jsonBody, err := json.Marshal(body)

	if err != nil {
		ch <- err
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		ch <- err
		return
	}

	if resp.StatusCode != 200 && len(body.SendTo) > 0 {
		// add foul
		fmt.Printf("failed to send new peer to %v, retraining with %v\n", peerUrl, body.SendTo[0])
		newBody := types.PeerPresentationBody{
			NewPeer: body.NewPeer,
			SendTo:  body.SendTo[1:],
		}
		wg.Add(1)
		p.SendNewPeer(newBody, body.SendTo[0], ch, wg)
		return
	}

	if resp.StatusCode != 200 && len(body.SendTo) == 0 {
		ch <- errors.New("Send to is empty")
		return
	}

	ch <- nil
}

func (p *PeerService) AllPeersToSend(excludeUrls []string) ([]models.Peer, error) {
	return p.repo.GetAll(excludeUrls)
}
