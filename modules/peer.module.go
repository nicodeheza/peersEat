package modules

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/validations"
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
	validate := validations.NewValidator(validator.New())
	peerControllers := controllers.NewPeerController(peerService, validate)
	return &PeerModule{
		peerCollection,
		peerRepository,
		peerService,
		peerControllers,
	}
}