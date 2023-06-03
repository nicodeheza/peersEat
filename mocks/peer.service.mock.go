package mocks

import (
	"errors"
	"sync"
	"time"

	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PeerServiceMock struct {
	Calls          map[string][][]interface{}
	InAreaPeerHave bool
}

func NewPeerServiceMock() *PeerServiceMock {
	return &PeerServiceMock{
		Calls: make(map[string][][]interface{}),
	}
}

func (p *PeerServiceMock) ClearCalls() {
	p.Calls = make(map[string][][]interface{})
}

func (p *PeerServiceMock) InitPeer() {
	//todo
	p.Calls["InitPeer"] = append(p.Calls["InitPeer"], []interface{}{nil})
	return
}

func (p *PeerServiceMock) EnqueueEvent(event types.Event) {}

func (p *PeerServiceMock) AllPeersToSend(excludeUrls []string) ([]models.Peer, error) {
	p.Calls["AllPeersToSend"] = append(p.Calls["AllPeersToSend"], []interface{}{excludeUrls})

	if excludeUrls[0] == "error" {
		return nil, errors.New("test error")
	}

	peer := models.Peer{
		Url: "http://tests.com",
	}

	return []models.Peer{peer, peer, peer, peer}, nil
}
func (p *PeerServiceMock) GetLocalPeer() (models.Peer, error) {
	if p.InAreaPeerHave {
		return models.Peer{
			Url:     "http://test.com",
			Center:  models.GeoCoords{Long: -34.577026, Lat: -58.466991},
			City:    "Buenos Aires",
			Country: "Argentina",
			InAreaPeers: []primitive.ObjectID{
				primitive.NewObjectIDFromTimestamp(time.Now()),
				{},
			},
		}, nil
	}
	return models.Peer{
		Url:     "http://test.com",
		Center:  models.GeoCoords{Long: -34.577026, Lat: -58.466991},
		City:    "Buenos Aires",
		Country: "Argentina",
		InAreaPeers: []primitive.ObjectID{
			primitive.NewObjectIDFromTimestamp(time.Now()),
			primitive.NewObjectIDFromTimestamp(time.Now()),
		},
	}, nil
}

func (p *PeerServiceMock) GetPeersUrlById(ids []primitive.ObjectID) ([]string, error) {
	for _, id := range ids {
		if id.IsZero() {
			return []string{"http://test.com", "http://have.com"}, nil
		}
	}

	return []string{"http://test.com", "http://test2.com"}, nil
}
func (p *PeerServiceMock) HaveRestaurant(restaurantQuery map[string]interface{}) (bool, error) {
	if restaurantQuery["name"] == "exist" {
		return true, nil
	}
	return false, nil
}

func (p *PeerServiceMock) PeerHaveRestaurant(peerUrl string, restaurantQuery map[string]interface{}, c chan<- types.PeerHaveRestaurantResp, wg *sync.WaitGroup) {
	defer wg.Done()
	if peerUrl == "http://have.com" {
		c <- types.PeerHaveRestaurantResp{
			Resp: true,
			Err:  nil,
		}
		return
	}
	c <- types.PeerHaveRestaurantResp{
		Resp: false,
		Err:  nil,
	}
}

func (p *PeerServiceMock) GetInDeliveryAreaPeers(peer models.Peer) ([]models.Peer, error) {
	return nil, nil
}

func (p *PeerServiceMock) GetNewDeliveryArea(peerCenter, restaurantCoord models.GeoCoords, restaurantDeliveryRadius float64) float64 {
	return 0
}

func (p *PeerServiceMock) UpdateDeliveryArea(peer models.Peer, newDeliveryRadius float64) error {
	return nil
}
