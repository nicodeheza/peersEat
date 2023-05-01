package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/mocks"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initTest() (*PeerController, *mocks.PeerServiceMock, *mocks.RestaurantServiceMock, *fiber.App) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	service := mocks.NewPeerServiceMock()
	validate := validations.NewValidator(validator.New())
	restaurantService := mocks.NewRestaurantServiceMock()
	geo := mocks.NewGeo()
	eventLoop := mocks.NewEventLoopMock()
	peerController := NewPeerController(service, validate, restaurantService, geo, eventLoop)
	app := fiber.New()

	return peerController, service, restaurantService, app
}

func TestPeerPresentation(t *testing.T) {
	controller, service, _, app := initTest()

	type Test struct {
		Title         string
		Body          types.PeerPresentationBody
		Status        int
		Json          string
		ServiceCallas map[string][][]interface{}
	}

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	wg1.Add(1)
	wg2.Add(2)

	sendTo := []string{"http://tests1.com", "http://tests2.com", "http://tests3.com", "http://tests4.com", "http://tests5.com", "http://tests6.com"}

	basePeer := models.Peer{
		Url:            "http://test.com",
		Center:         models.GeoCoords{Long: 1, Lat: 1},
		City:           "test",
		Country:        "test",
		DeliveryRadius: 3,
	}
	tests := []Test{
		{
			Title: "send invalid peer",
			Body: types.PeerPresentationBody{
				NewPeer: models.Peer{
					Url: "test",
				},
				SendTo: nil,
			},
			Status: 400,
			Json:   "[map[FailedField:Peer.Url Tag:url Value:] map[FailedField:Peer.Center.Long Tag:required Value:] map[FailedField:Peer.Center.Lat Tag:required Value:] map[FailedField:Peer.City Tag:required Value:] map[FailedField:Peer.Country Tag:required Value:]]",
		},
		{
			Title: "success with not send to",
			Body: types.PeerPresentationBody{
				NewPeer: basePeer,
				SendTo:  nil,
			},
			Status: 200,
			ServiceCallas: map[string][][]interface{}{
				"AddNewPeer": {{basePeer}},
			},
		},
		{
			Title: "success with send to",
			Body: types.PeerPresentationBody{
				NewPeer: basePeer,
				SendTo:  sendTo,
			},
			Status: 200,
			ServiceCallas: map[string][][]interface{}{
				"AddNewPeer": {{basePeer}},
				"GetSendMap": {{sendTo, make(map[string][]string)}},
				"SendNewPeer": {{
					types.PeerPresentationBody{
						NewPeer: basePeer,
						SendTo:  []string{"http://tests2.com", "http://tests3.com"},
					},
					"http://tests1.com",
					make(chan error),
					&wg1,
				},
					{
						types.PeerPresentationBody{
							NewPeer: basePeer,
							SendTo:  []string{"http://tests5.com", "http://tests6.com"},
						},
						"http://tests4.com",
						make(chan error),
						&wg2,
					},
				},
			},
		},
	}

	app.Post("/", controller.PeerPresentation)

	for _, test := range tests {

		body, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, 1)
		if err != nil {
			t.Fatal()
		}

		var b interface{}
		json.NewDecoder(resp.Body).Decode(&b)

		if test.Status != resp.StatusCode {
			t.Errorf("%s\n expecting status %d but got %d",
				test.Title, test.Status, resp.StatusCode)
		}

		if test.Json != "" && test.Json != fmt.Sprintf("%v", b) {
			t.Errorf("%s\n incorrect response body\n expecting: %s\n got: %s",
				test.Title, test.Json, fmt.Sprintf("%v", b))
		}

		if test.ServiceCallas != nil {
			for k, v := range test.ServiceCallas {
				expect := v
				got := service.Calls[k]

				for i := range expect {
					for j, ele := range expect[i] {
						if reflect.ValueOf(ele).Kind() == reflect.Ptr {
							expect[i][j] = "ptr"
						}
						if reflect.ValueOf(ele).Kind() == reflect.Chan {
							expect[i][j] = "chan"
						}
					}
				}

				for i := range got {
					for j, ele := range got[i] {
						if reflect.ValueOf(ele).Kind() == reflect.Ptr {
							got[i][j] = "ptr"
						}
						if reflect.ValueOf(ele).Kind() == reflect.Chan {
							got[i][j] = "chan"
						}
					}
				}

				if !reflect.DeepEqual(expect, got) {
					tmp := expect[0]
					expect[0] = expect[1]
					expect[1] = tmp
					if !reflect.DeepEqual(expect, got) {
						t.Errorf("%s\n\n incorrect service %s call\n\n expecting: %v\n\n got: %v\n\n",
							test.Title, k, expect, got)

					}
				}
			}
		}

		service.ClearCalls()
	}
}

func TestSendAllPeers(t *testing.T) {
	controller, service, _, app := initTest()

	type Test struct {
		Title         string
		Query         string
		Status        int
		Json          string
		ServiceCallas map[string][][]interface{}
	}

	queryError := "error,http://test2.com"
	query := "http://test1.com,http://test2.com"

	tests := []Test{
		{
			Title:  "return error",
			Query:  queryError,
			Status: 500,
			Json:   "map[message:test error]",
			ServiceCallas: map[string][][]interface{}{
				"AllPeersToSend": {{strings.Split(queryError, ",")}},
			},
		},
		{
			Title:  "return success",
			Query:  query,
			Status: 200,
			Json:   "[map[Center:map[Lat:0 Long:0] id:000000000000000000000000 url:http://tests.com] map[Center:map[Lat:0 Long:0] id:000000000000000000000000 url:http://tests.com] map[Center:map[Lat:0 Long:0] id:000000000000000000000000 url:http://tests.com] map[Center:map[Lat:0 Long:0] id:000000000000000000000000 url:http://tests.com]]",
			ServiceCallas: map[string][][]interface{}{
				"AllPeersToSend": {{strings.Split(query, ",")}},
			},
		},
	}

	app.Get("/", controller.SendAllPeers)
	for _, test := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/?excludes=%s", test.Query), nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, 1)
		if err != nil {
			t.Fatal()
		}

		var b interface{}

		json.NewDecoder(resp.Body).Decode(&b)

		if resp.StatusCode != test.Status {
			t.Errorf("%s\n\n incorrect status code\n\n expected: %d\n got: %d\n\n",
				test.Title, test.Status, resp.StatusCode)
		}

		bodyString := fmt.Sprintf("%v", b)

		if test.Json != bodyString {
			t.Errorf("%s\n\n incorrect body\n\n expected: %s\n\n got: %s\n\n",
				test.Title, test.Json, bodyString)
		}

		if !reflect.DeepEqual(test.ServiceCallas, service.Calls) {
			t.Errorf("%s\n\n incorrect service call\n\n expected: %v\n\n got: %v\n\n",
				test.Title, test.ServiceCallas, service.Calls)
		}

		service.ClearCalls()
	}
}

func TestHaveRestaurant(t *testing.T) {
	controller, service, _, app := initTest()

	type Test struct {
		Title  string
		Query  string
		Status int
		Json   string
	}

	tests := []Test{
		{
			Title:  "not exist",
			Query:  "name=test&address=test&city=test&country=test",
			Status: 200,
			Json:   "map[result:false]",
		},
		{
			Title:  "exist",
			Query:  "name=exist&address=test&city=test&country=test",
			Status: 200,
			Json:   "map[result:true]",
		},
	}

	app.Get("/", controller.HaveRestaurant)
	for _, test := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/?%s", test.Query), nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, 1)
		if err != nil {
			t.Fatal()
		}

		var b interface{}

		json.NewDecoder(resp.Body).Decode(&b)

		if resp.StatusCode != test.Status {
			t.Errorf("%s\n\n incorrect status code\n\n expected: %d\n got: %d\n\n",
				test.Title, test.Status, resp.StatusCode)
		}

		bodyString := fmt.Sprintf("%v", b)

		if test.Json != bodyString {
			t.Errorf("%s\n\n incorrect body\n\n expected: %s\n\n got: %s\n\n",
				test.Title, test.Json, bodyString)
		}

		service.ClearCalls()
	}
}

func TestAddNewRestaurant(t *testing.T) {
	controller, service, restaurantService, app := initTest()

	type Test struct {
		Title      string
		Body       models.Restaurant
		Status     int
		Json       string
		WasAdded   bool
		InAreaHave bool
	}

	tests := []Test{
		{
			Title: "Out of influence area",
			Body: models.Restaurant{
				Name:    "out",
				Address: "testAddress",
				City:    "testCity",
				Country: "testCountry",
			},
			Status: 400,
			Json:   "map[message:restaurant out of area]",
		},
		{
			Title: "Send existent restaurant",
			Body: models.Restaurant{
				Name:    "test",
				Address: "testAddress",
				City:    "testCity",
				Country: "testCountry",
			},
			Status:     400,
			Json:       "map[message:restaurant already exists]",
			InAreaHave: true,
		},
		{
			Title: "add it successfully",
			Body: models.Restaurant{
				Name:    "test",
				Address: "testAddress",
				City:    "testCity",
				Country: "testCountry",
			},
			Status:   200,
			Json:     "map[newRestaurant:map[Address:testAddress City:testCity Coord:map[Lat:1 Long:1] Country:testCountry Name:test id:000000000000000000000000 menu:map[Sections:<nil>] password:testHash rate:map[Stars:0 Votes:0] userName:testUsername] tempPassword:testPassword]",
			WasAdded: true,
		},
	}

	app.Post("/", controller.AddNewRestaurant)
	for _, test := range tests {

		service.InAreaPeerHave = test.InAreaHave

		body, err := json.Marshal(test.Body)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, 1)
		if err != nil {
			t.Fatal()
		}

		var b interface{}

		json.NewDecoder(resp.Body).Decode(&b)

		if resp.StatusCode != test.Status {
			t.Errorf("%s\n\n incorrect status code\n\n expected: %d\n got: %d\n\n",
				test.Title, test.Status, resp.StatusCode)
		}

		bodyString := fmt.Sprintf("%v", b)

		if test.Json != bodyString {
			t.Errorf("%s\n\n incorrect body\n\n expected: %s\n\n got: %s\n\n",
				test.Title, test.Json, bodyString)
		}

		if test.WasAdded {
			expectRestaurant := models.Restaurant{
				Id:              primitive.ObjectID{},
				Name:            "test",
				Address:         "testAddress",
				City:            "testCity",
				Country:         "testCountry",
				Coord:           models.GeoCoords{Long: 1, Lat: 1},
				Password:        "testHash",
				UserName:        "testUsername",
				IsFinalPassword: false,
			}
			savedRestaurant := restaurantService.Calls["AddNewRestaurant"][0][0]
			if !reflect.DeepEqual(expectRestaurant, savedRestaurant) {
				t.Errorf("%s\n\n incorrect restaurant saved\n Expecting: %v\n Got: %v\n",
					test.Title, expectRestaurant, savedRestaurant)
			}
		} else {
			if restaurantService.Calls["AddNewRestaurant"] != nil {
				t.Error("Restaurant was saved when wasn't expected")
			}
		}

		service.ClearCalls()
	}
}
