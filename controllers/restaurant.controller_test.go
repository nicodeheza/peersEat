package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/mocks"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
)

func initTestRestaurant() (*RestaurantController, *mocks.RestaurantServiceMock, *fiber.App) {
	service := mocks.NewRestaurantServiceMock()
	peerServices := mocks.NewPeerServiceMock()
	validations := validations.NewValidator(validator.New())
	controller := NewRestaurantController(service, peerServices, validations)
	app := fiber.New()

	return controller, service, app
}

func TestUpdatePassword(t *testing.T) {
	controller, _, app := initTestRestaurant()

	app.Patch("/", controller.UpdatePassword)

	sBody := types.UpdateRestaurantPassword{
		NewPassword: "testPassword",
		Id:          "5f9a8a5c7c9d440000a9a8c7",
	}
	body, err := json.Marshal(sBody)
	if err != nil {
		t.Error(err.Error())
	}

	req := httptest.NewRequest("PATCH", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 1)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	var b map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&b)

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if b["message"] != "success" {
		t.Errorf("Expected message success, got %s", b["message"])
	}
}
