package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantController struct {
	Service     services.RestaurantServiceI
	peerService services.PeerServiceI
	validators  validations.ValidateI
}

type RestaurantControllerI interface {
	UpdatePassword(c *fiber.Ctx) error
	RetuneOk(c *fiber.Ctx) error
	UpdateRestaurantData(c *fiber.Ctx) error
}

func NewRestaurantController(service services.RestaurantServiceI, peerService services.PeerServiceI, validators validations.ValidateI) *RestaurantController {
	return &RestaurantController{service, peerService, validators}
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
	err = r.Service.UpdateRestaurantUsernameAndPassword(id, body.NewPassword, body.NewUsername)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

func (r *RestaurantController) RetuneOk(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

// update restaurant data
func (r *RestaurantController) UpdateRestaurantData(c *fiber.Ctx) error {
	restaurantData := new(types.RestaurantData)
	if err := c.BodyParser(&restaurantData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	deliveryUpdated, coords, radius, err := r.Service.UpdateData(*restaurantData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if deliveryUpdated {
		selfPeer, err := r.peerService.GetLocalPeer()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		newDeliveryArea := r.peerService.GetNewDeliveryArea(selfPeer.Center, coords, radius)

		err = r.peerService.UpdateDeliveryArea(selfPeer, newDeliveryArea)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return c.SendStatus(200)
}

// update menu
