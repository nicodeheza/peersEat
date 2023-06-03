package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/mocks"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/services/validations"
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
	peerController := NewPeerController(service, validate, restaurantService, geo)
	app := fiber.New()

	return peerController, service, restaurantService, app
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
