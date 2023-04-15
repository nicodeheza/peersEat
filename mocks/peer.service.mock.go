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

func (p *PeerServiceMock) AddNewPeer(newPeer models.Peer) error {
	p.Calls["AddNewPeer"] = append(p.Calls["AddNewPeer"], []interface{}{newPeer})
	if newPeer.Url == "http://error.com" {
		return errors.New("test error")
	}
	return nil
}

func (p *PeerServiceMock) GetSendMap(urls []string, sendMap map[string][]string) {
	gotMap := make(map[string][]string)
	for k, v := range sendMap {
		gotMap[k] = v
	}
	p.Calls["GetSendMap"] = append(p.Calls["GetSendMap"], []interface{}{urls, gotMap})

	sendMap["http://tests1.com"] = []string{"http://tests2.com", "http://tests3.com"}
	sendMap["http://tests4.com"] = []string{"http://tests5.com", "http://tests6.com"}
}

func (p *PeerServiceMock) GetNewSendMap(excludes []string, sendMap map[string][]string) error {
	gotMap := make(map[string][]string)
	for k, v := range sendMap {
		gotMap[k] = v
	}
	p.Calls["GetNewSendMap"] = append(p.Calls["GetNewSendMap"], []interface{}{excludes, gotMap})

	if excludes[0] == "error" {
		return errors.New("test error")
	}
	urls := []string{"http://tests1.com", "http://tests2.com", "http://tests3.com", "http://tests4.com", "http://tests5.com", "http://tests6.com"}
	p.GetSendMap(urls, sendMap)
	return nil
}

func (p *PeerServiceMock) SendNewPeer(body types.PeerPresentationBody, peerUrl string, ch chan<- error, wg *sync.WaitGroup) {
	p.Calls["SendNewPeer"] = append(p.Calls["SendNewPeer"], []interface{}{body, peerUrl, ch, wg})

	defer wg.Done()
	if peerUrl == "error" {
		ch <- errors.New("test error")
		return
	}
	ch <- nil
	return
}

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
