package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/middleware"
)

func peerRoutes(app *fiber.App, controllers controllers.PeerControllerI, authMiddleware middleware.AuthMiddlewareI) {
	peerGroup := app.Group("/peer")

	peerGroup.Post("/present", controllers.PeerPresentation)
	peerGroup.Get("/all", controllers.SendAllPeers)
	peerGroup.Get("/restaurant/have", controllers.HaveRestaurant)
	peerGroup.Post("/restaurant", authMiddleware.OnlyPeerOwner, controllers.AddNewRestaurant)
}
