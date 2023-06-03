package modules

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/events"
	"github.com/nicodeheza/peersEat/middleware"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/utils"
)

func initRepositories() *Repositories {
	restaurantCollection := models.GetRestaurantColl("peersEatDB")
	restaurantRepository := repositories.NewRestaurantRepository(restaurantCollection)

	peerCollection := models.GetPeerColl("peersEatDB")
	peerRepository := repositories.NewPeerRepository(peerCollection)

	return &Repositories{Peer: peerRepository, Restaurant: restaurantRepository}
}

func initServices(repos *Repositories, authHelpers *utils.AuthHelpers, eventLoop *events.EventLoop, geo *geo.GeoService) *Services {
	restaurant := services.NewRestaurantService(repos.Restaurant, authHelpers, geo)
	peer := services.NewPeerService(repos.Peer, geo, repos.Restaurant, eventLoop)

	return &Services{peer, restaurant}
}

func initControllers(services *Services, validate *validations.Validate, geo *geo.GeoService) *Controllers {
	restaurant := controllers.NewRestaurantController(services.restaurant, services.peer, validate)
	peer := controllers.NewPeerController(services.peer, validate, services.restaurant, geo)
	return &Controllers{peer, restaurant}
}

func InitApp() *Application {

	geo := geo.NewGeo()
	validate := validations.NewValidator(validator.New())
	authHelpers := utils.NewAuthHelper()

	repos := initRepositories()
	eventHandlers := events.NewEventHandlers(repos.Peer, validate, geo)
	eventLoop := events.InitEventLoop(eventHandlers)
	services := initServices(repos, authHelpers, eventLoop, geo)
	controllers := initControllers(services, validate, geo)

	restaurantModule := &RestaurantModule{repos.Restaurant, services.restaurant, controllers.restaurant}
	peerModule := &PeerModule{repos.Peer, services.peer, controllers.peer}
	authMiddleware := middleware.InitAuthMiddleware(restaurantModule.Service)

	return &Application{peerModule, restaurantModule, authMiddleware}
}
