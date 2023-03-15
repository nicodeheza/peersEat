package services

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/mocks"
	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initTest()(*PeerService, *mocks.PeerRepositoryMock){
	err := godotenv.Load("../.env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}
	repo := mocks.NewPeerRepository()
	geo := mocks.NewGeo()
	return NewPeerService(repo,geo), repo
}

func TestInitPeer(t *testing.T){
	service, repo := initTest()
	service.InitPeer()
	savePeer := repo.InsertCalls[0]

	if savePeer.Url != os.Getenv("HOST"){
		t.Errorf("incorrect host:\n expected: %s\n received: %s", os.Getenv("HOST"),savePeer.Url)
	}

	centerSrt := fmt.Sprintf("%f,%f", savePeer.Center.Long, savePeer.Center.Lat)

	if centerSrt != os.Getenv("CENTER"){
		t.Errorf("incorrect center:\n expected: %s\n received: %s", os.Getenv("CENTER"), centerSrt)
	}
	if savePeer.City != os.Getenv("CITY"){
		t.Errorf("incorrect city:\n expected: %s\n received: %s", os.Getenv("CITY"), savePeer.City)
	}
	if savePeer.Country != os.Getenv("COUNTRY"){
		t.Errorf("incorrect country:\n expected: %s\n received: %s", os.Getenv("COUNTRY"), savePeer.Country)
	}
}

func TestAddNewPeer(t *testing.T){
	service, repo := initTest()

	basePeer := models.Peer{
		Url: "test",
		Center: models.Center{Long: 0,Lat: 0},
		City: os.Getenv("CITY"),
		Country: os.Getenv("COUNTRY"),
		InfluenceRadius: 0,
		DeliveryRadius: 0,
	}
	difCity:= basePeer
	difCity.City= "test"
	difCountry:= basePeer
	difCountry.Country ="test"
	difCityAndCountry:= basePeer
	difCityAndCountry.City= "test"
	difCityAndCountry.Country= "test"
	inInfluence := basePeer
	inInfluence.InfluenceRadius= 1
	inDelivery := basePeer
	inDelivery.DeliveryRadius =1
	inBoth:= inInfluence
	inBoth.DeliveryRadius=1

	
	type test struct{
		Peer models.Peer
		Expect mocks.ExpectUpdate
	}

	selfPeer,_:= repo.GetSelf()
	repo.CleatCalls()

	selfInInfluence :=selfPeer
	selfInInfluence.InAreaPeers= append(selfInInfluence.InAreaPeers, primitive.NilObjectID)
	selfInDelivery :=selfPeer
	selfInDelivery.InDeliveryAreaPeers= append(selfInDelivery.InDeliveryAreaPeers, primitive.NilObjectID)
	selfInBoth :=selfPeer
	selfInBoth.InAreaPeers= append(selfInBoth.InAreaPeers, primitive.NilObjectID)
	selfInBoth.InDeliveryAreaPeers= append(selfInBoth.InDeliveryAreaPeers, primitive.NilObjectID)

	tests := []test{
		{
			Peer: basePeer,
		}, 
		{
			Peer: difCity,
		},
		{
			Peer: difCountry,
		},
		{
			Peer: difCityAndCountry,
		},
		{
			Peer: inInfluence,
			Expect: mocks.ExpectUpdate{Peer: selfInInfluence, Fields: []string{"InAreaPeers"}},
		},
		{
			Peer: inDelivery,
			Expect: mocks.ExpectUpdate{Peer: selfInDelivery, Fields: []string{"InDeliveryAreaPeers"}},
		},
		{
			Peer: inBoth,
			Expect: mocks.ExpectUpdate{Peer: selfInBoth, Fields: []string{"InAreaPeers","InDeliveryAreaPeers"}},
		},
	}

	for _, test := range tests{
		service.AddNewPeer(test.Peer)
		insertCall:= repo.InsertCalls[0]
		if !reflect.DeepEqual(insertCall, test.Peer){
			t.Errorf("incorrect insertion.\n Expect: %v\n get: %v", test.Peer, insertCall)
		}

		updateCalls := repo.UpdateCalls

		if reflect.DeepEqual((mocks.ExpectUpdate{}),test.Expect) && len(updateCalls) > 0{
			t.Error("update should not be called.")
		}


		if len(updateCalls) >0 && !reflect.DeepEqual(updateCalls[0], test.Expect){
			t.Errorf("update should be called with args\n %v\n, but had been caller with args\n %v", test.Expect, updateCalls[0])
		}
		repo.CleatCalls()
	}
}

func TestGetSendMap(t *testing.T){
	service, _ := initTest()

	type GetSendMapArgs struct{
		Urls []string
		SendMap map[string][]string
	}

	map1:= make(map[string][]string)
	map1["test1"] = nil
	map2 := make(map[string][]string)
	map2["test1"]=[]string{"test2"}
	map2["test3"]=[]string{"test4"}

	tests:=[]struct{
		args GetSendMapArgs
		result map[string][]string
	}{
		{
		args: GetSendMapArgs{
			Urls: []string{},
			SendMap: make(map[string][]string),
		},
		result: make(map[string][]string),
		},
		{
		args: GetSendMapArgs{
			Urls: []string{"test1"},
			SendMap: make(map[string][]string),
		},
		result: map1,
		},
		{
			args: GetSendMapArgs{
				Urls: []string{"test1", "test2", "test3" , "test4"},
				SendMap: make(map[string][]string),
			},
			result: map2,
		},
	}

	for _, test := range tests{
		service.GetSendMap(test.args.Urls, test.args.SendMap)

		if ! reflect.DeepEqual(test.args.SendMap, test.result){
			t.Errorf("unexpected result.\n expected: %v\n got: %v", test.result, test.args.SendMap)
		}
	}

}