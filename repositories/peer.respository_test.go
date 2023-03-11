package repositories

import (
	"context"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
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
		t.Errorf("document insertion failed with err: %v", err)
	}

	result, err:= peerRepository.GetById(id)
	if err!=nil{
		t.Errorf("document get by id failed with err: %v", err)
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
	newPeer4 := models.Peer{
		Url: "http://tests4.com",
		Center: models.Center{Long: 44.0,Lat: 44.0},
		City: "test city4",
		Country: "test country4",
		InfluenceRadius: 5,
		DeliveryRadius: 6,
	}

	newPeers :=[]*models.Peer{&newPeer1, &newPeer2, &newPeer3, &newPeer4}

	for _, peer := range newPeers{
		res, err := peerRepository.Insert(*peer)
		if err != nil{
			t.Errorf("document insertion failed with err: %v", err)
		}
		peer.Id= res
	}

	res, err := peerRepository.GetAll([]string{newPeer3.Url, newPeer4.Url})

	if err != nil{
		t.Errorf("document get all failed with err: %v", err)
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

func TestGetSelf(t *testing.T){
	coll, server:= initDb()
	defer server.Stop(context.Background())
	
	peerRepository:= PeerRepository{coll}

	centerSlice := strings.Split(os.Getenv("CENTER"), ",")
	long, err := strconv.ParseFloat(centerSlice[0], 64)
	if err !=nil{t.Errorf("fail to parse float with err:\n%v", err)}
	lat, err:= strconv.ParseFloat(centerSlice[1], 64)
	if err !=nil{t.Errorf("fail to parse float with err:\n%v", err)}


	selfPeer := models.Peer{
		Url: os.Getenv("HOST"),
		Center: models.Center{Long:long ,Lat: lat},
		City: os.Getenv("CITY"),
		Country: os.Getenv("COUNTRY"),
		InfluenceRadius: 1,
		DeliveryRadius: 2,
	}

	id, err:= peerRepository.Insert(selfPeer)
	if err != nil{t.Errorf("fail to insert peer with err:\n%v", err)}

	selfPeer.Id= id

	result, err := peerRepository.GetSelf()
	if err != nil{t.Errorf("fail to get peer with err:\n%v", err)}

	if !reflect.DeepEqual(selfPeer, result){
		t.Errorf("get and target peers are not equals:\n get: %v\n target: %v ", result,selfPeer)
	}
}

func TestUpdate(t *testing.T){
	coll, server:= initDb()
	defer server.Stop(context.Background())
	
	peerRepository:= PeerRepository{coll}

	peer := models.Peer{
		Url: "http://tests.com",
		Center: models.Center{Long: 99.0,Lat: 99.0},
		City: "test city",
		Country: "test country",
		InfluenceRadius: 1,
		DeliveryRadius: 2,
	}

	id, err := peerRepository.Insert(peer)
	if err != nil{t.Errorf("fail to insert peer with err:\n%v", err)}

	peer.Id = id

	expectedError := peerRepository.Update(peer, []string{"badFiled"})
	if expectedError == nil{
		t.Error("Update should throw and error")
	}

	peer.Center= models.Center{Long: 11.0,Lat: 11.0}
	peer.City = "cityUpdated"

	err = peerRepository.Update(peer, []string{"Center", "City"})
	if err != nil{t.Errorf("update failed with error:\n%v", err)}

	result, err := peerRepository.GetById(peer.Id)
	if err != nil{t.Errorf("get by id failed with error:\n%v", err)}

	if !reflect.DeepEqual(result, peer){
		t.Errorf("expect peer and receive peer are not equal.\n receive: %v\n expect: %v\n", result, peer)
	}
}

func TestGetAllUrls(t *testing.T){
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
	newPeer4 := models.Peer{
		Url: "http://tests4.com",
		Center: models.Center{Long: 44.0,Lat: 44.0},
		City: "test city4",
		Country: "test country4",
		InfluenceRadius: 5,
		DeliveryRadius: 6,
	}

	newPeers :=[]*models.Peer{&newPeer1, &newPeer2, &newPeer3, &newPeer4}

	for _, peer := range newPeers{
		res, err := peerRepository.Insert(*peer)
		if err != nil{
			t.Errorf("document insertion failed with err: %v", err)
		}
		peer.Id= res
	}

	res, err := peerRepository.GetAllUrls([]string{newPeer3.Url,newPeer4.Url})

	if err != nil{
		t.Errorf("document get all urls failed with err: %v", err)
	}

	if len(res) !=2{
		t.Errorf("get all return urls the incorrect amount of urls:\n%v", res)
	}

	for i, url := range res{
		if url != *&newPeers[i].Url{
			t.Errorf("elements are not equal:\n %v\n %v", url, newPeers[i].Url)
		}
	}
}