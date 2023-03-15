package mocks

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpectUpdate struct{
	Peer models.Peer
	Fields []string
}

type PeerRepositoryMock struct{
	InsertCalls []models.Peer
	GetByIdCalls []primitive.ObjectID
	GetAllCalls [][] string
	UpdateCalls []ExpectUpdate
	GetAllUrlsCalls [][]string
}

func NewPeerRepository() *PeerRepositoryMock{
	return &PeerRepositoryMock{}
}

func (p *PeerRepositoryMock) CleatCalls(){
	p.InsertCalls = nil
	p.GetByIdCalls = nil
	p.UpdateCalls = nil
	p.GetAllUrlsCalls = nil
	p.GetAllCalls= nil
}

func (p *PeerRepositoryMock) Insert(peer models.Peer) (id primitive.ObjectID ,  err error) {
	p.InsertCalls= append(p.InsertCalls, peer)
	if peer.Url == "error"{
		return primitive.ObjectID{}, errors.New("test error")
	}
	return primitive.NilObjectID, nil
}

func (p *PeerRepositoryMock) GetById(id  primitive.ObjectID)(models.Peer, error){
	p.GetByIdCalls = append(p.GetByIdCalls, id)
	return models.Peer{
		Id: id,
		Url: "http://tests.com",
		Center: models.Center{Long: 11, Lat: 11},
		City: "test city",
		Country: "test country",
		InfluenceRadius: 2,
		DeliveryRadius: 3,
	}, nil
}

func (p *PeerRepositoryMock) GetAll(excludesUrls []string) ([]models.Peer, error) {
	p.GetAllCalls= append(p.GetAllCalls, excludesUrls)
	peer := models.Peer{
		Id: primitive.NilObjectID,
		Url: "http://tests.com",
		Center: models.Center{Long: 11, Lat: 11},
		City: "test city",
		Country: "test country",
		InfluenceRadius: 2,
		DeliveryRadius: 3,
	}

	return []models.Peer{peer, peer, peer}, nil
}

func (p *PeerRepositoryMock) GetSelf() (models.Peer, error) {
	center := strings.Split(os.Getenv("CENTER"), ",") 
	long, err:= strconv.ParseFloat(center[0], 64)
	lat, err:= strconv.ParseFloat(center[1], 64)
	if err != nil{
		return models.Peer{}, err
	}
	return models.Peer{
		Id:primitive.NilObjectID,
		Url: os.Getenv("HOST"),
		Center: models.Center{Long: long, Lat: lat},
		City: os.Getenv("CITY"),
		Country: os.Getenv("COUNTRY"),
		InfluenceRadius: 2,
		DeliveryRadius: 4,
	}, nil
} 

func (p *PeerRepositoryMock) Update(peer models.Peer, fields []string ) error {
	p.UpdateCalls= append(p.UpdateCalls, ExpectUpdate{peer, fields})
	return nil
}


func (p *PeerRepositoryMock) GetAllUrls(excludes []string) ([]string, error){
	
	p.GetAllUrlsCalls= append(p.GetAllUrlsCalls, excludes)

	url:="http://test.com"

	return []string{url, url, url}, nil
}