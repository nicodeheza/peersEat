package services

import (
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantService struct {
	repo        repositories.RestaurantRepositoryI
	authHelpers utils.AuthHelpersI
	geo         geo.GeoServiceI
}

func NewRestaurantService(repository repositories.RestaurantRepositoryI, authHelpers utils.AuthHelpersI, geo geo.GeoServiceI) *RestaurantService {
	return &RestaurantService{repository, authHelpers, geo}
}

func (r *RestaurantService) CompleteRestaurantInitialData(newRestaurant *models.Restaurant) error {
	coords, err := r.geo.GetAddressCoords(newRestaurant.Address, newRestaurant.City, newRestaurant.Name)
	if err != nil {
		return err
	}

	newRestaurant.Coord = coords

	newPassword := r.authHelpers.GetRandomPassword(20)
	newHash, err := r.authHelpers.HashPasswords(newPassword)
	if err != nil {
		return err
	}
	newUserName, err := r.authHelpers.GetRandomWords(3)
	if err != nil {
		return err
	}

	newRestaurant.Password = newHash
	newRestaurant.UserName = newUserName
	newRestaurant.IsFinalPassword = false

	return nil
}

func (r *RestaurantService) AddNewRestaurant(newRestaurant models.Restaurant) (primitive.ObjectID, error) {
	id, err := r.repo.Insert(newRestaurant)
	if err != nil {
		return id, err
	}
	return id, nil
}
