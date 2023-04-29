package services

import (
	"errors"

	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/types"
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
	UpdateRestaurantUsernameAndPassword(id primitive.ObjectID, newPassword string, newUserNames string) error
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

func (r *RestaurantService) UpdateRestaurantUsernameAndPassword(id primitive.ObjectID, newPassword string, newUserNames string) error {
	hash, err := r.authHelpers.HashPasswords(newPassword)
	if err != nil {
		return err
	}
	return r.repo.Update(id, map[string]interface{}{"password": hash, "userName": newUserNames, "isFinalPassword": true})
}

func (r *RestaurantService) Authenticate(password, userName string) (bool, string, error) {
	restaurant, err := r.repo.FindOne(map[string]interface{}{"userName": userName})
	if err != nil {
		return false, "", err
	}
	if !restaurant.IsFinalPassword {
		return false, "", errors.New("Update your username and password before login")
	}
	res := r.authHelpers.CheckPassword(password, restaurant.Password)
	id := restaurant.Id.Hex()

	return res, id, nil
}

func (r *RestaurantService) UpdateData(data types.RestaurantData) error {
	updates := make(map[string]interface{})
	updates["name"] = data.Name
	updates["ImageUrl"] = data.ImageUrl
	updates["openTime"] = data.OpenTime
	updates["closeTime"] = data.CloseTime
	updates["phone"] = data.Phone
	updates["deliveryCost"] = data.DeliveryCost
	updates["isDeliveryFixCost"] = data.IsDeliveryFixCost
	updates["minDeliveryTime"] = data.MinDeliveryTime
	updates["maxDeliveryTime"] = data.MaxDeliveryTime
	updates["deliveryRadius"] = data.DeliveryRadius

	id, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return err
	}
	err = r.repo.Update(id, updates)
	if err != nil {
		return err
	}

	return nil
}

// update menu
