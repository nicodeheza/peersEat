package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nicodeheza/peersEat/config"
	"github.com/nicodeheza/peersEat/middleware"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/modules"
	"github.com/nicodeheza/peersEat/routes"
)

func main() {
	config.LoadEnv()
	app := fiber.New()
	app.Use(logger.New())

	app.Use(middleware.InitSessionMiddleware().Middleware)

	config.ConnectDB(os.Getenv("MONGO_URI"))
	models.InitModels("peersEatDB")
	appModule := modules.InitApp()
	appModule.Peer.Service.InitPeer()

	routes.Register(app, appModule)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!!")
	})

	port := os.Getenv("PORT")
	if "" == port {
		port = "3001"
	}
	port = fmt.Sprintf(":%v", port)
	app.Listen(port)
}
