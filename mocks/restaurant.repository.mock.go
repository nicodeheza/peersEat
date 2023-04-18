package mocks

import (
	"errors"

	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantRepositoryMock struct {
}

func NewRestaurantRepositoryMock() *RestaurantRepositoryMock {
	return &RestaurantRepositoryMock{}
}

func (r *RestaurantRepositoryMock) Insert(restaurant models.Restaurant) (id primitive.ObjectID, err error) {
	return primitive.ObjectID{}, errors.New("not implemented")
}

func (r *RestaurantRepositoryMock) FindOne(query map[string]interface{}) (models.Restaurant, error) {
	return models.Restaurant{}, errors.New("not implemented")
}
func (r *RestaurantRepositoryMock) Update(id primitive.ObjectID, updates map[string]interface{}) error {
	return nil
}
