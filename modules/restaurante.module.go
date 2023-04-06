package modules

import (
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantModule struct {
	Collection *mongo.Collection
	Repository repositories.RestaurantRepositoryI
	Service    services.RestaurantServiceI
}

func NewRestaurantModule(authHelpers utils.AuthHelpersI, geo geo.GeoServiceI) *RestaurantModule {
	collection := models.GetRestaurantColl("peersEatDB")
	repository := repositories.NewRestaurantRepository(collection)
	service := services.NewRestaurantService(repository, authHelpers, geo)

	return &RestaurantModule{
		collection,
		repository,
		service,
	}
}
