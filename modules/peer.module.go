package modules

import (
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/events"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
	"go.mongodb.org/mongo-driver/mongo"
)

type PeerModule struct {
	Collection  *mongo.Collection
	Repository  repositories.PeerRepositoryI
	Service     services.PeerServiceI
	Controllers controllers.PeerControllerI
}

func newPeerModule(validate validations.ValidateI,
	geo geo.GeoServiceI,
	restaurants services.RestaurantServiceI,
	restaurantRepo repositories.RestaurantRepositoryI,
) *PeerModule {
	peerCollection := models.GetPeerColl("peersEatDB")
	peerRepository := repositories.NewPeerRepository(peerCollection)
	peerService := services.NewPeerService(peerRepository, geo, restaurantRepo)
	eventHandlers := events.NewEventHandlers(peerRepository, validate, geo)
	eventLoop := events.InitEventLoop(eventHandlers)
	peerControllers := controllers.NewPeerController(peerService, validate, restaurants, geo, eventLoop)
	return &PeerModule{
		peerCollection,
		peerRepository,
		peerService,
		peerControllers,
	}
}
