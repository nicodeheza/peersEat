package modules

import (
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type PeerModule struct {
	Collection  *mongo.Collection
	Repository  repositories.PeerRepositoryI
	Service     services.PeerServiceI
	Controllers controllers.PeerControllerI
}

func newPeerModule() *PeerModule{
	peerCollection := models.GetPeerColl("peersEatDB")
	peerRepository := repositories.NewPeerRepository(peerCollection)
	peerService := services.NewPeerService(peerRepository)
	peerControllers := controllers.NewPeerController(peerService)
	return &PeerModule{
		peerCollection,
		peerRepository,
		peerService,
		peerControllers,
	}
}