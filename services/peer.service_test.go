package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/mocks"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initTest() (*PeerService, *mocks.PeerRepositoryMock) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	os.Setenv("INITIAL_PEER", "http://test.com")
	repo := mocks.NewPeerRepository()
	geo := mocks.NewGeo()
	restaurantRepository := mocks.NewRestaurantRepositoryMock()
	return NewPeerService(repo, geo, restaurantRepository), repo
}

func TestInitPeer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	service, repo := initTest()
	defer repo.ClearCalls()

	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/peer/present", os.Getenv("INITIAL_PEER")),
		httpmock.NewStringResponder(200, ``))

	allPeers := []models.Peer{}

	for i := 0; i < 10; i++ {
		allPeers = append(allPeers, models.Peer{
			Id:             primitive.NewObjectIDFromTimestamp(time.Now()),
			Url:            fmt.Sprintf("http://test%d.com", i),
			Center:         models.GeoCoords{Long: float64(i), Lat: float64(i)},
			City:           fmt.Sprintf("City%d", i),
			Country:        fmt.Sprintf("Country%d", i),
			DeliveryRadius: float64(i + 1),
		})
	}

	allPeersJson, err := json.Marshal(allPeers)
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/peer/all", os.Getenv("INITIAL_PEER")),
		httpmock.NewStringResponder(200, string(allPeersJson)))

	service.InitPeer()
	savePeer := repo.InsertCalls[0]

	if savePeer.Url != os.Getenv("HOST") {
		t.Errorf("incorrect host:\n expected: %s\n received: %s", os.Getenv("HOST"), savePeer.Url)
	}

	centerSrt := fmt.Sprintf("%f,%f", savePeer.Center.Long, savePeer.Center.Lat)

	if centerSrt != os.Getenv("CENTER") {
		t.Errorf("incorrect center:\n expected: %s\n received: %s", os.Getenv("CENTER"), centerSrt)
	}
	if savePeer.City != os.Getenv("CITY") {
		t.Errorf("incorrect city:\n expected: %s\n received: %s", os.Getenv("CITY"), savePeer.City)
	}
	if savePeer.Country != os.Getenv("COUNTRY") {
		t.Errorf("incorrect country:\n expected: %s\n received: %s", os.Getenv("COUNTRY"), savePeer.Country)
	}
	if !reflect.DeepEqual(repo.InsertManyCalls[0], allPeers) {
		t.Errorf("incorrect insertMany args:\n expected: %v\n received: %v", allPeers, repo.InsertManyCalls[0])
	}

}

func TestAddNewPeer(t *testing.T) {
	service, repo := initTest()

	basePeer := models.Peer{
		Url:            "test",
		Center:         models.GeoCoords{Long: 0, Lat: 0},
		City:           os.Getenv("CITY"),
		Country:        os.Getenv("COUNTRY"),
		DeliveryRadius: 0,
	}
	difCity := basePeer
	difCity.City = "test"
	difCountry := basePeer
	difCountry.Country = "test"
	difCityAndCountry := basePeer
	difCityAndCountry.City = "test"
	difCityAndCountry.Country = "test"
	inInfluence := basePeer
	inDelivery := basePeer
	inDelivery.DeliveryRadius = 1
	inBoth := inInfluence
	inBoth.DeliveryRadius = 2

	type test struct {
		Peer   models.Peer
		Expect mocks.ExpectUpdate
	}

	selfPeer, _ := repo.GetSelf()
	repo.ClearCalls()

	selfInInfluence := selfPeer
	selfInInfluence.InAreaPeers = append(selfInInfluence.InAreaPeers, primitive.NilObjectID)
	selfInDelivery := selfPeer
	selfInDelivery.InDeliveryAreaPeers = append(selfInDelivery.InDeliveryAreaPeers, primitive.NilObjectID)
	selfInBoth := selfPeer
	selfInBoth.InAreaPeers = append(selfInBoth.InAreaPeers, primitive.NilObjectID)
	selfInBoth.InDeliveryAreaPeers = append(selfInBoth.InDeliveryAreaPeers, primitive.NilObjectID)
	selfInBoth.DeliveryRadius = 2

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
			Peer:   inInfluence,
			Expect: mocks.ExpectUpdate{Peer: selfInInfluence, Fields: []string{"InAreaPeers"}},
		},
		{
			Peer:   inDelivery,
			Expect: mocks.ExpectUpdate{Peer: selfInDelivery, Fields: []string{"InDeliveryAreaPeers"}},
		},
		{
			Peer:   inBoth,
			Expect: mocks.ExpectUpdate{Peer: selfInBoth, Fields: []string{"InAreaPeers", "InDeliveryAreaPeers"}},
		},
	}

	for _, test := range tests {
		service.AddNewPeer(test.Peer)
		insertCall := repo.InsertCalls[0]
		if !reflect.DeepEqual(insertCall, test.Peer) {
			t.Errorf("incorrect insertion.\n Expect: %v\n get: %v", test.Peer, insertCall)
		}

		updateCalls := repo.UpdateCalls

		if reflect.DeepEqual((mocks.ExpectUpdate{}), test.Expect) && len(updateCalls) > 0 {
			t.Error("update should not be called.")
		}

		if len(updateCalls) > 0 && !reflect.DeepEqual(updateCalls[0], test.Expect) {
			t.Errorf("update should be called with args\n %v\n, but had been caller with args\n %v", test.Expect, updateCalls[0])
		}
		repo.ClearCalls()
	}
}

func TestGetSendMap(t *testing.T) {
	service, _ := initTest()

	type GetSendMapArgs struct {
		Urls    []string
		SendMap map[string][]string
	}

	map1 := make(map[string][]string)
	map1["test1"] = nil
	map2 := make(map[string][]string)
	map2["test1"] = []string{"test2"}
	map2["test3"] = []string{"test4"}

	tests := []struct {
		args   GetSendMapArgs
		result map[string][]string
	}{
		{
			args: GetSendMapArgs{
				Urls:    []string{},
				SendMap: make(map[string][]string),
			},
			result: make(map[string][]string),
		},
		{
			args: GetSendMapArgs{
				Urls:    []string{"test1"},
				SendMap: make(map[string][]string),
			},
			result: map1,
		},
		{
			args: GetSendMapArgs{
				Urls:    []string{"test1", "test2", "test3", "test4"},
				SendMap: make(map[string][]string),
			},
			result: map2,
		},
	}

	for _, test := range tests {
		service.GetSendMap(test.args.Urls, test.args.SendMap)

		if !reflect.DeepEqual(test.args.SendMap, test.result) {
			t.Errorf("unexpected result.\n expected: %v\n got: %v", test.result, test.args.SendMap)
		}
	}

}

func TestGetNewSendMap(t *testing.T) {
	service, repo := initTest()

	sendMap := make(map[string][]string)
	excludes := []string{"http://tests.com"}

	service.GetNewSendMap(excludes, sendMap)

	if len(repo.GetAllUrlsCalls) > 1 {
		t.Error("GetNewSendMap called more than one time")
	}
	if !reflect.DeepEqual(repo.GetAllUrlsCalls[0], excludes) {
		t.Errorf("GetAllUrls call with unexpected argument.\n expected: %v\n got: %v", excludes, repo.GetAllUrlsCalls[0])
	}

	expectMap := make(map[string][]string)
	expectMap["test1"] = []string{"test2"}
	expectMap["test3"] = []string{"test4"}

	if !reflect.DeepEqual(expectMap, sendMap) {
		t.Errorf("unexpected send map.\n expected: %v\n got: %v", expectMap, sendMap)
	}
}

func TestSendNewPeer(t *testing.T) {
	service, _ := initTest()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://test.com/peer/present",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})
	httpmock.RegisterResponder("POST", "http://error.com/peer/present",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), nil
		})

	httpmock.RegisterResponder("POST", "http://error2.com/peer/present",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), errors.New("test error")
		})

	newPeer := models.Peer{
		Url: "test",
	}

	tests := []struct {
		Title string
		Body  types.PeerPresentationBody
		Url   string
		Err   error
	}{
		{
			"test1",
			types.PeerPresentationBody{
				NewPeer: newPeer,
				SendTo: []string{
					"http://test.com", "http://test2.com",
				},
			},
			"http://error.com",
			nil,
		},
		{
			"test2",
			types.PeerPresentationBody{
				NewPeer: newPeer,
				SendTo:  []string{},
			},
			"http://error.com",
			errors.New("Send to is empty"),
		},
		{
			"test3",
			types.PeerPresentationBody{
				NewPeer: newPeer,
				SendTo:  []string{"https://test.com", "http://test2.com"},
			},
			"http://error2.com",
			errors.New(`Post "http://error2.com/peer/present": test error`),
		},
	}

	var wg sync.WaitGroup
	ch := make(chan error)

	for _, test := range tests {
		wg.Add(1)
		go service.SendNewPeer(test.Body, test.Url, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make([]error, 0)
	for err := range ch {
		result = append(result, err)
	}

	for i, test := range tests {
		err := test.Err
		if err == nil && result[i] != nil {
			t.Errorf("Not expecting error on test: %s\n expecting: %s\n got: %s\n", test.Title, err, result[i].Error())
		}

		if err != nil && result[i] == nil {
			t.Errorf("Expecting error on test: %s\n expecting: %s\n got: %s\n", test.Title, err.Error(), result[i])
		}

		if err != nil && result[i] != nil && err.Error() != result[i].Error() {
			t.Errorf("Different error expected on test: %s\n expecting: %s\n got: %s\n", test.Title, err.Error(), result[i].Error())
		}
	}

}

func TestAllPeerToSend(t *testing.T) {
	service, repo := initTest()

	arg := []string{"http://fake.com"}
	expect, _ := repo.GetAll(arg)
	repo.ClearCalls()

	result, _ := service.AllPeersToSend(arg)

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("expecting result: %v\ngot: %v\n", expect, result)
	}

	if !reflect.DeepEqual(repo.GetAllCalls[0], arg) {
		t.Errorf("expecting arg: %v\ngot: %v\n", arg, repo.GetAllCalls[0])
	}

}

func TestPeerHaveRestaurant(t *testing.T) {
	//todo
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	service, repo := initTest()
	defer repo.ClearCalls()

	httpmock.RegisterResponder("GET", "http://test.com/peer/restaurant/have",
		httpmock.NewStringResponder(200, `{"result": false}`))

	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan types.PeerHaveRestaurantResp)
	query := map[string]interface{}{
		"name":    "test",
		"address": "test",
	}

	go service.PeerHaveRestaurant("http://test.com", query, c, &wg)

	go func() {
		wg.Wait()
		close(c)
	}()

	res := <-c
	if res.Resp {
		t.Error("expecting false but got true")
	}
	if res.Err != nil {
		t.Errorf("%s", res.Err.Error())
	}
}
