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

type RestaurantServiceI interface {
	CompleteRestaurantInitialData(newRestaurant *models.Restaurant) (string, error)
	AddNewRestaurant(newRestaurant models.Restaurant) (primitive.ObjectID, error)
	UpdateRestaurantPassword(id primitive.ObjectID, newPassword string) error
	Authenticate(password, userName string) (bool, string, error)
}

func NewRestaurantService(repository repositories.RestaurantRepositoryI, authHelpers utils.AuthHelpersI, geo geo.GeoServiceI) *RestaurantService {
	return &RestaurantService{repository, authHelpers, geo}
}

func (r *RestaurantService) CompleteRestaurantInitialData(newRestaurant *models.Restaurant) (string, error) {
	coords, err := r.geo.GetAddressCoords(newRestaurant.Address, newRestaurant.City, newRestaurant.Country)
	if err != nil {
		return "", err
	}

	newRestaurant.Coord = coords

	newPassword := r.authHelpers.GetRandomPassword(20)
	newHash, err := r.authHelpers.HashPasswords(newPassword)
	if err != nil {
		return "", err
	}
	newUserName, err := r.authHelpers.GetRandomWords(3)
	if err != nil {
		return "", err
	}

	newRestaurant.Password = newHash
	newRestaurant.UserName = newUserName
	newRestaurant.IsFinalPassword = false

	return newPassword, nil
}

func (r *RestaurantService) AddNewRestaurant(newRestaurant models.Restaurant) (primitive.ObjectID, error) {
	id, err := r.repo.Insert(newRestaurant)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *RestaurantService) UpdateRestaurantPassword(id primitive.ObjectID, newPassword string) error {
	hash, err := r.authHelpers.HashPasswords(newPassword)
	if err != nil {
		return err
	}
	return r.repo.Update(id, map[string]interface{}{"password": hash})
}

func (r *RestaurantService) Authenticate(password, userName string) (bool, string, error) {
	restaurant, err := r.repo.FindOne(map[string]interface{}{"userName": userName})
	if err != nil {
		return false, "", err
	}
	res := r.authHelpers.CheckPassword(password, restaurant.Password)
	id := restaurant.Id.Hex()

	return res, id, nil
}
