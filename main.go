package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nicodeheza/peersEat/config"
	"github.com/nicodeheza/peersEat/models"
	peer_service "github.com/nicodeheza/peersEat/services/peerService"
)

func main() {
	config.LoadEnv()
    app := fiber.New()
	app.Use(logger.New())

	config.ConnectDB()
	models.InitModels()
	peer_service.InitPeer()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!!!!")
    })

	port:= os.Getenv("PORT")
	if "" == port{
		port= "3001"
	}
	port = fmt.Sprintf(":%v", port) 
    app.Listen(port)
}