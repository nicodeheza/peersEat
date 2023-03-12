package modules

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
)

type Application struct {
	Peer *PeerModule
}

func InitApp() *Application {

	geo := geo.NewGeo()
	validate := validations.NewValidator(validator.New())
	peerModule := newPeerModule(validate, geo)

	return &Application{peerModule}
}