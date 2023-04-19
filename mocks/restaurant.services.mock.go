package mocks

import (
	"errors"

	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantServiceMock struct {
	Calls map[string][][]interface{}
}

func NewRestaurantServiceMock() *RestaurantServiceMock {
	calls := make(map[string][][]interface{})
	return &RestaurantServiceMock{Calls: calls}
}

func (r *RestaurantServiceMock) CompleteRestaurantInitialData(newRestaurant *models.Restaurant) (string, error) {
	if newRestaurant.Name == "out" {
		newRestaurant.Coord = models.GeoCoords{Long: 0, Lat: 0}
	} else {
		newRestaurant.Coord = models.GeoCoords{Long: 1, Lat: 1}
	}
	newRestaurant.Password = "testHash"
	newRestaurant.UserName = "testUsername"
	newRestaurant.IsFinalPassword = false
	return "testPassword", nil
}

func (r *RestaurantServiceMock) AddNewRestaurant(newRestaurant models.Restaurant) (primitive.ObjectID, error) {
	r.Calls["AddNewRestaurant"] = append(r.Calls["AddNewRestaurant"], []interface{}{newRestaurant})
	return primitive.ObjectID{}, nil
}
func (r *RestaurantServiceMock) UpdateRestaurantPassword(id primitive.ObjectID, newPassword string) error {
	if newPassword == "error" {
		return errors.New("test error")
	}
	return nil
}
