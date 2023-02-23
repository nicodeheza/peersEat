package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nicodeheza/peersEat/config"
)

func main() {
    app := fiber.New()
	app.Use(logger.New())

	config.ConnectDB()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!!!!")
    })

	port:= config.GetEnv("PORT")
	if "" == port{
		port= "3001"
	}
	port = fmt.Sprintf(":%v", port) 
    app.Listen(port)
}