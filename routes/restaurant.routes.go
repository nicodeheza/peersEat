package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/controllers"
	"github.com/nicodeheza/peersEat/middleware"
)

func RestaurantRoutes(app *fiber.App, controllers controllers.RestaurantControllerI, authMiddleware middleware.AuthMiddlewareI) {
	restaurantGroup := app.Group("/restaurant")
	restaurantGroup.Patch("/password", authMiddleware.Protect, controllers.UpdatePassword)
}
