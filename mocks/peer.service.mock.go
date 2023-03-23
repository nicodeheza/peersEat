package mocks

import (
	"errors"
	"sync"

	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
)

type PeerServiceMock struct {
	Calls map[string][][]interface{}
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
