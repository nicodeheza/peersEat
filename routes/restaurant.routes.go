package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/controllers"
)

func RestaurantRoutes(app *fiber.App, controllers controllers.RestaurantControllerI) {
	restaurantGroup := app.Group("/restaurant")
	restaurantGroup.Patch("/password", controllers.UpdatePassword)
}
