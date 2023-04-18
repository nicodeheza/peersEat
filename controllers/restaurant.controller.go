package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantController struct {
	Service services.RestaurantServiceI
}

type RestaurantControllerI interface {
	UpdatePassword(c *fiber.Ctx) error
}

func NewRestaurantController(service services.RestaurantServiceI) *RestaurantController {
	return &RestaurantController{Service: service}
}

func (r *RestaurantController) UpdatePassword(c *fiber.Ctx) error {
	body := new(types.UpdateRestaurantPassword)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	id, err := primitive.ObjectIDFromHex(body.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	err = r.Service.UpdateRestaurantPassword(id, body.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}
