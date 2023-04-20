package modules

import (
	"github.com/nicodeheza/peersEat/controllers"
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
	Controller controllers.RestaurantControllerI
}

func NewRestaurantModule(authHelpers utils.AuthHelpersI, geo geo.GeoServiceI) *RestaurantModule {
	collection := models.GetRestaurantColl("peersEatDB")
	repository := repositories.NewRestaurantRepository(collection)
	service := services.NewRestaurantService(repository, authHelpers, geo)
	controller := controllers.NewRestaurantController(service)

	return &RestaurantModule{
		collection,
		repository,
		service,
		controller,
	}
}
