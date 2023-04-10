package modules

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/utils"
)

type Application struct {
	Peer *PeerModule
}

func InitApp() *Application {

	geo := geo.NewGeo()
	validate := validations.NewValidator(validator.New())
	authHelpers := utils.NewAuthHelper()
	restaurantModule := NewRestaurantModule(authHelpers, geo)
	peerModule := newPeerModule(validate, geo, restaurantModule.Service, restaurantModule.Repository)

	return &Application{peerModule}
}
