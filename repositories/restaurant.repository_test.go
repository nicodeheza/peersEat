package repositories

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/joho/godotenv"
	"github.com/nicodeheza/peersEat/config"
	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func initRestaurantDb() (*mongo.Collection, *mim.Server) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server, err := mim.StartWithOptions(context.Background(), "6.0.0", mim.WithPort(27017))
	if err != nil {
		log.Fatal("Error creating in memory db")
	}
	config.ConnectDB(server.URI())
	dbName := "peersEatDBTest"
	models.InitModels(dbName)
	return models.GetRestaurantColl(dbName), server
}

func TestInsertAndFindOne(t *testing.T) {
	coll, server := initRestaurantDb()
	defer server.Stop(context.Background())

	rr := RestaurantRepository{coll}

	newRestaurants := []models.Restaurant{
		{
			Name:     "test",
			Address:  "testAddress",
			City:     "testCity",
			Country:  "testCountry",
			UserName: "test1",
		},
		{
			Name:     "test2",
			Address:  "testAddress2",
			City:     "testCity2",
			Country:  "testCountry2",
			UserName: "test2",
		},
	}

	for i, restaurant := range newRestaurants {
		id, err := rr.Insert(restaurant)
		fmt.Println(id)
		if err != nil {
			t.Errorf("fail to insert restaurant with error: %v", err)
		}
		if id.IsZero() {
			t.Error("Id not created")
		}
		newRestaurants[i].Id = id
	}

	query1 := map[string]interface{}{
		"name":    newRestaurants[0].Name,
		"address": newRestaurants[0].Address,
	}
	rest, err := rr.FindOne(query1)
	if err != nil {
		t.Errorf("fail to find restaurant with error: %v", err)
	}
	if !reflect.DeepEqual(rest, newRestaurants[0]) {
		t.Errorf("incorrect restaurant found.\n expected: %v\n got: %v\n", newRestaurants[0], rest)
	}

	query2 := map[string]interface{}{
		"name":    "test3",
		"address": "testAddress3",
	}
	rest, err = rr.FindOne(query2)
	if err.Error() != "mongo: no documents in result" {
		t.Errorf("expecting error: mongo: no documents in result but got: %v", err)
	}
}

func TestUpdateRestaurant(t *testing.T) {
	coll, server := initRestaurantDb()
	defer server.Stop(context.Background())

	rr := RestaurantRepository{coll}
	id, err := rr.Insert(models.Restaurant{
		Name:            "test",
		Address:         "testAddress",
		City:            "testCity",
		Country:         "testCountry",
		UserName:        "testUsername",
		Password:        "testPassword",
		IsFinalPassword: false,
	})
	if err != nil {
		t.Errorf("fail to insert restaurant with error: %v", err)
	}

	updates := map[string]interface{}{
		"name":            "testUpdated",
		"address":         "testAddressUpdated",
		"city":            "testCityUpdated",
		"country":         "testCountryUpdated",
		"userName":        "testUsernameUpdated",
		"password":        "testPasswordUpdated",
		"isFinalPassword": true,
	}

	err = rr.Update(id, updates)
	if err != nil {
		t.Errorf("fail to update restaurant with error: %v", err)
	}
	res, err := rr.FindOne(map[string]interface{}{"_id": id})
	if err != nil {
		t.Errorf("fail to find restaurant with error: %v", err)
	}
	fmt.Println(">>>>>", res)
	if res.Name != updates["name"] {
		t.Errorf("incorrect restaurant name.\n expected: %v\n got: %v\n", updates["name"], res.Name)
	}
	if res.Address != updates["address"] {
		t.Errorf("incorrect restaurant address.\n expected: %v\n got: %v\n", updates["address"], res.Address)
	}
	if res.City != updates["city"] {
		t.Errorf("incorrect restaurant city.\n expected: %v\n got: %v\n", updates["city"], res.City)
	}
	if res.UserName != updates["userName"] {
		t.Errorf("incorrect restaurant userName.\n expected: %v\n got: %v\n", updates["userName"], res.UserName)
	}
	if res.Password != updates["password"] {
		t.Errorf("incorrect restaurant password.\n expected: %v\n got: %v\n", updates["password"], res.Password)
	}
	if res.IsFinalPassword != updates["isFinalPassword"] {
		t.Errorf("incorrect restaurant isFinalPassword.\n expected: %v\n got: %v\n", updates["isFinalPassword"], res.IsFinalPassword)
	}
}
