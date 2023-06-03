package modules

import (
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/middleware"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services"
)

type Application struct {
	Peer           *PeerModule
	Restaurant     *RestaurantModule
	AuthMiddleware *middleware.AuthMiddleware
}

type Repositories struct {
	Peer       repositories.PeerRepositoryI
	Restaurant repositories.RestaurantRepositoryI
}

type Services struct {
	peer       services.PeerServiceI
	restaurant services.RestaurantServiceI
}

type Controllers struct {
	peer       controllers.PeerControllerI
	restaurant controllers.RestaurantControllerI
}

type RestaurantModule struct {
	Repository repositories.RestaurantRepositoryI
	Service    services.RestaurantServiceI
	Controller controllers.RestaurantControllerI
}

type PeerModule struct {
	Repository  repositories.PeerRepositoryI
	Service     services.PeerServiceI
	Controllers controllers.PeerControllerI
}
