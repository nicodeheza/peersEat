package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/mocks"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
)

func initTest()(*PeerController, *mocks.PeerServiceMock, *fiber.App){
	err := godotenv.Load("../.env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	service:= mocks.NewPeerServiceMock()
	validate:= validations.NewValidator(validator.New())
	peerController:= NewPeerController(service, validate)
	app:= fiber.New()

	return peerController, service, app
}

func TestPeerPresentation(t *testing.T) {
	controller, service, app := initTest() 

	type Test struct{
		Title string
		Body types.PeerPresentationBody
		Status int
		Json string
		ServiceCallas map[string][][]interface{}
	}

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	wg1.Add(1)
	wg2.Add(2)

	sendTo:=[] string{"http://tests1.com","http://tests2.com","http://tests3.com","http://tests4.com","http://tests5.com","http://tests6.com",}

	basePeer:= models.Peer{
		Url: "http://test.com",
		Center: models.Center{Long: 1,Lat: 1},
		City:"test",
		Country: "test",
		InfluenceRadius: 2,
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
			Json: "[map[FailedField:Peer.Url Tag:url Value:] map[FailedField:Peer.Center.Long Tag:required Value:] map[FailedField:Peer.Center.Lat Tag:required Value:] map[FailedField:Peer.City Tag:required Value:] map[FailedField:Peer.Country Tag:required Value:]]",
		},
		{
			Title: "success with not send to",
			Body: types.PeerPresentationBody{
				NewPeer: basePeer,
				SendTo: nil,
			},
			Status: 200,
			ServiceCallas: map[string][][]interface{}{
				"AddNewPeer": {{basePeer}},
				"GetNewSendMap": {{[]string{basePeer.Url, os.Getenv("HOST")}, make(map[string][]string)}},
				"SendNewPeer": {{
					types.PeerPresentationBody{
						NewPeer: basePeer,
						SendTo: []string{"http://tests2.com","http://tests3.com"},
					},
					"http://tests1.com",
					make(chan error),
					&wg1,
				},
				{
					types.PeerPresentationBody{
						NewPeer: basePeer,
						SendTo: []string{"http://tests5.com","http://tests6.com"},
					},
					"http://tests4.com",
					make(chan error),
					&wg2,
				},
			},
			},
		},
		{
			Title: "success with send to",
			Body: types.PeerPresentationBody{
				NewPeer: basePeer,
				SendTo: sendTo,
			},
			Status: 200,
			ServiceCallas: map[string][][]interface{}{
				"AddNewPeer": {{basePeer}},
				"GetSendMap": {{sendTo, make(map[string][]string)}},
				"SendNewPeer": {{
					types.PeerPresentationBody{
						NewPeer: basePeer,
						SendTo: []string{"http://tests2.com","http://tests3.com"},
					},
					"http://tests1.com",
					make(chan error),
					&wg1,
				},
				{
					types.PeerPresentationBody{
						NewPeer: basePeer,
						SendTo: []string{"http://tests5.com","http://tests6.com"},
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

	for _,test := range tests{

		body, err := json.Marshal(test.Body)
		if err != nil{
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/",  bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req,1)
		if err != nil{
			t.Fatal()
		}

		var b interface{}
		json.NewDecoder(resp.Body).Decode(&b)
		
		if test.Status != resp.StatusCode{
			t.Errorf("%s\n expecting status %d but got %d", 
			test.Title, test.Status, resp.StatusCode)
		}

		if test.Json !="" && test.Json != fmt.Sprintf("%v",b){
			t.Errorf("%s\n incorrect response body\n expecting: %s\n got: %s",
					test.Title, test.Json, fmt.Sprintf("%v",b))
		}

		if test.ServiceCallas != nil{
			for k,v:=range test.ServiceCallas{
				expect := v
				got := service.Calls[k]

				for i:= range expect{
					for j, ele:= range expect[i]{
						if reflect.ValueOf(ele).Kind()==reflect.Ptr{
							expect[i][j]= "ptr"
						}
						if reflect.ValueOf(ele).Kind()==reflect.Chan{
							expect[i][j]= "chan"
						}
					}
				} 

				for i:= range got{
					for j, ele:= range got[i]{
						if reflect.ValueOf(ele).Kind()==reflect.Ptr{
							got[i][j]= "ptr"
						}
						if reflect.ValueOf(ele).Kind()==reflect.Chan{
							got[i][j]= "chan"
						}
					}
				} 

				if !reflect.DeepEqual(expect, got){
					tmp:= expect[0]
					expect[0]= expect[1]
					expect[1]=tmp
					if !reflect.DeepEqual(expect, got){
						t.Errorf("%s\n\n incorrect service %s call\n\n expecting: %v\n\n got: %v\n\n",
								test.Title, k ,expect, got)
						
					}
				}
			}
		}

		service.ClearCalls()
	}
}
