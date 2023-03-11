package repositories

import (
	"context"
	"log"
	"reflect"
	"testing"

	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/config"
	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func initDb() (*mongo.Collection, *mim.Server) {
	err := godotenv.Load("../.env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}
	server, err := mim.StartWithOptions(context.Background(),"6.0.0", mim.WithPort(27017))
	if err != nil{
		log.Fatal("Error creating in memory db")
	}
	config.ConnectDB(server.URI())
	dbName:="peersEatDBTest"
	models.InitModels(dbName)
	return models.GetPeerColl(dbName), server
}

func TestInsertAndFindById(t *testing.T){
	coll, server:= initDb()
	defer server.Stop(context.Background())

	peerRepository:= PeerRepository{coll}

	newPeer := models.Peer{
		Url: "http://tests.com",
		Center: models.Center{Long: 99.0,Lat: 99.0},
		City: "test city",
		Country: "test country",
		InfluenceRadius: 1,
		DeliveryRadius: 2,
	}

	id, err := peerRepository.Insert(newPeer)

	if err!=nil{
		t.Errorf("documente insertion failed with err: %v", err)
	}

	result, err:= peerRepository.GetById(id)
	if err!=nil{
		t.Errorf("documente get by id failed with err: %v", err)
	}

	newPeer.Id= result.Id
	if !reflect.DeepEqual(newPeer, result){
		t.Errorf("elements are not equal:\n %v\n %v", newPeer, result)
	}  

}

func TestGetAll(t *testing.T){
	coll, server:= initDb()
	defer server.Stop(context.Background())
	
	peerRepository:= PeerRepository{coll}

	newPeer1 := models.Peer{
		Url: "http://tests1.com",
		Center: models.Center{Long: 11.0,Lat: 11.0},
		City: "test city1",
		Country: "test country1",
		InfluenceRadius: 1,
		DeliveryRadius: 2,
	}
	newPeer2 := models.Peer{
		Url: "http://tests2.com",
		Center: models.Center{Long: 22.0,Lat: 22.0},
		City: "test city2",
		Country: "test country2",
		InfluenceRadius: 2,
		DeliveryRadius: 3,
	}
	newPeer3 := models.Peer{
		Url: "http://tests3.com",
		Center: models.Center{Long: 33.0,Lat: 33.0},
		City: "test city3",
		Country: "test country3",
		InfluenceRadius: 4,
		DeliveryRadius: 5,
	}

	newPeers :=[]*models.Peer{&newPeer1, &newPeer2, &newPeer3}

	for _, peer := range newPeers{
		res, err := peerRepository.Insert(*peer)
		if err != nil{
			t.Errorf("documente insertion failed with err: %v", err)
		}
		peer.Id= res
	}

	res, err := peerRepository.GetAll([]string{newPeer3.Url})

	if err != nil{
		t.Errorf("documente get all failed with err: %v", err)
	}

	if len(res) !=2{
		t.Errorf("get all return the incorrect amount of peers:\n%v", res)
	}

	for i, peer := range res{
		if !reflect.DeepEqual(peer, *newPeers[i]){
			t.Errorf("elements are not equal:\n %v\n %v", peer, newPeers[i])
		}
	}
}
