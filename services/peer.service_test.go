package services

import (
	"encoding/json"
	"fmt"
	"log"
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
	eventsLoop := mocks.NewEventLoopMock()
	return NewPeerService(repo, geo, restaurantRepository, eventsLoop), repo
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
