package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PeerServiceI interface {
	InitPeer()
	AddNewPeer(newPeer models.Peer) error
	GetSendMap(urls []string, sendMap map[string][]string)
	GetNewSendMap(excludes []string, sendMap map[string][]string) error
	SendNewPeer(body types.PeerPresentationBody, peerUrl string, ch chan<- error, wg *sync.WaitGroup)
	AllPeersToSend(excludeUrls []string) ([]models.Peer, error)
	GetLocalPeer() (models.Peer, error)
	GetPeersUrlById(ids []primitive.ObjectID) ([]string, error)
	HaveRestaurant(restaurantQuery map[string]interface{}) (bool, error)
}

type PeerService struct {
	repo           repositories.PeerRepositoryI
	geo            geo.GeoServiceI
	restaurantRepo repositories.RestaurantRepositoryI
}

func NewPeerService(repository repositories.PeerRepositoryI, geo geo.GeoServiceI, restaurantRepo repositories.RestaurantRepositoryI) *PeerService {
	return &PeerService{repository, geo, restaurantRepo}
}

func (p *PeerService) InitPeer() {
	defer fmt.Println("Peer installed successfully")
	centerSrt := os.Getenv("CENTER")
	centerSlice := strings.Split(centerSrt, ",")
	long, _ := strconv.ParseFloat(centerSlice[0], 64)
	lat, _ := strconv.ParseFloat(centerSlice[1], 64)

	selfPeer := models.Peer{
		Url:     os.Getenv("HOST"),
		Center:  models.GeoCoords{Long: long, Lat: lat},
		City:    os.Getenv("CITY"),
		Country: os.Getenv("COUNTRY"),
	}

	p.repo.Insert(selfPeer)

	initialPeer := os.Getenv("INITIAL_PEER")

	if initialPeer != "" {

		resp, err := http.Get(fmt.Sprintf("%s/peer/all?excludes=%s", initialPeer, selfPeer.Url))
		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err)
			fmt.Println(resp.StatusCode)
			log.Fatal("bad request")
		}

		newPeers := make([]models.Peer, 0)
		decoder := json.NewDecoder(resp.Body)

		err = decoder.Decode(&newPeers)
		if err != nil {
			log.Fatal("fail to decode")
		}

		_, err = p.repo.InsertMany(newPeers)
		if err != nil {
			log.Fatal("fail to inset new peers")
		}

		sendTo, err := p.repo.GetAllUrls([]string{selfPeer.Url, initialPeer})
		if err != nil {
			log.Fatal("fail to get all peers")
		}

		bodyMap := map[string]interface{}{
			"newPeer": map[string]interface{}{
				"url": selfPeer.Url,
				"center": map[string]interface{}{
					"long": selfPeer.Center.Long,
					"lat":  selfPeer.Center.Lat,
				},
				"city":    selfPeer.City,
				"Country": selfPeer.Country,
			},
			"sendTo": sendTo,
		}

		postBody, err := json.Marshal(bodyMap)
		if err != nil {
			log.Fatal("Marshal error")
		}
		resp, err = http.Post(fmt.Sprintf("%s/peer/present", initialPeer),
			"application/json", bytes.NewBuffer(postBody))
		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err)
			fmt.Println(resp.StatusCode)
			log.Fatal("bad request")
		}

	}
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
	if p.geo.AreInfluenceAreasOverlaying(selfPeer, newPeer) {
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

func (p *PeerService) GetLocalPeer() (models.Peer, error) {
	return p.repo.GetSelf()
}

func (p *PeerService) GetPeersUrlById(ids []primitive.ObjectID) ([]string, error) {
	return p.repo.FindUrlsByIds(ids)
}

func (p *PeerService) HaveRestaurant(restaurantQuery map[string]interface{}) (bool, error) {
	_, err := p.restaurantRepo.FindOne(restaurantQuery)

	if err.Error() == "mongo: no documents in result" {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

// validate restaurant (check multiple reqs)
