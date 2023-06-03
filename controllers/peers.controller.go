package controllers

import (
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
)

type PeerControllerI interface {
	SendAllPeers(c *fiber.Ctx) error
	HaveRestaurant(c *fiber.Ctx) error
	AddNewRestaurant(c *fiber.Ctx) error
	EventReceiver(c *fiber.Ctx) error
}

type PeerController struct {
	service     services.PeerServiceI
	validate    validations.ValidateI
	restaurants services.RestaurantServiceI
	geo         geo.GeoServiceI
}

func NewPeerController(service services.PeerServiceI, validate validations.ValidateI, restaurants services.RestaurantServiceI, geo geo.GeoServiceI) *PeerController {
	return &PeerController{service, validate, restaurants, geo}
}

func (p *PeerController) EventReceiver(c *fiber.Ctx) error {
	body := new(types.Event)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	errors := p.validate.ValidateEvent(*body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	p.service.EnqueueEvent(*body)
	return c.SendStatus(fiber.StatusOK)
}

func (p *PeerController) SendAllPeers(c *fiber.Ctx) error {
	query := new(types.SendAllPeerQuery)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	peers, err := p.service.AllPeersToSend(query.Excludes)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(peers)
}

func (p *PeerController) AddNewRestaurant(c *fiber.Ctx) error {
	newRestaurant := models.Restaurant{}
	err := c.BodyParser(&newRestaurant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	errors := p.validate.ValidateRestaurant(newRestaurant)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	password, err := p.restaurants.CompleteRestaurantInitialData(&newRestaurant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	self, err := p.service.GetLocalPeer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	isInInfluenceArea := p.geo.IsInInfluenceArea(self.Center, newRestaurant.Coord)

	if !isInInfluenceArea {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "restaurant out of area"})
	}

	inAreaUrl, err := p.service.GetPeersUrlById(self.InAreaPeers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	query := make(map[string]interface{})
	if newRestaurant.Name != "" {
		query["name"] = newRestaurant.Name
	}
	if newRestaurant.Address != "" {
		query["address"] = newRestaurant.Address
	}
	if newRestaurant.City != "" {
		query["city"] = newRestaurant.City
	}
	if newRestaurant.Country != "" {
		query["country"] = newRestaurant.Country
	}

	ch := make(chan types.PeerHaveRestaurantResp)
	var wg sync.WaitGroup

	for _, url := range inAreaUrl {
		wg.Add(1)
		go p.service.PeerHaveRestaurant(url, query, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var isValid bool
	var errs []string
	for resp := range ch {
		if resp.Err != nil {
			errs = append(errs, resp.Err.Error())
		}
		if resp.Resp {
			isValid = false
			break
		}
		isValid = true
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": errs})
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "restaurant already exists"})
	}

	id, err := p.restaurants.AddNewRestaurant(newRestaurant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	newRestaurant.Id = id

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"newRestaurant": newRestaurant,
		"tempPassword":  password,
	})
}

func (p *PeerController) HaveRestaurant(c *fiber.Ctx) error {
	restaurantQuery := make(map[string]interface{})
	querySrt := string(c.Request().URI().QueryString())
	for _, pairSrt := range strings.Split(querySrt, "&") {
		pair := strings.Split(pairSrt, "=")
		restaurantQuery[pair[0]] = pair[1]
	}

	have, err := p.service.HaveRestaurant(restaurantQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": have})
}
