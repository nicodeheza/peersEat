package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/controllers"
)

func peerRoutes(app *fiber.App, controllers controllers.PeerControllerI) {
	peerGroup := app.Group("/peer")

	peerGroup.Post("/present", controllers.PeerPresentation)
	peerGroup.Get("/all", controllers.SendAllPeers)
}
