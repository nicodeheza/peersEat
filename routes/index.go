package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/modules"
)

func Register(app *fiber.App, appModule *modules.Application) {

	peerRoutes(app, appModule.Peer.Controllers, appModule.AuthMiddleware)
	RestaurantRoutes(app, appModule.Restaurant.Controller, appModule.AuthMiddleware)
}
