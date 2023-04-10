package controllers

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
)

type PeerControllerI interface {
	PeerPresentation(c *fiber.Ctx) error
	SendAllPeers(c *fiber.Ctx) error
	HaveRestaurant(c *fiber.Ctx) error
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

type peerPresentationBody struct {
	NewPeer models.Peer
	SendTo  []string
}

func (p *PeerController) PeerPresentation(c *fiber.Ctx) error {
	body := new(types.PeerPresentationBody)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	newPeer := body.NewPeer

	errors := p.validate.ValidatePeer(newPeer)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	p.service.AddNewPeer(newPeer)

	if body.SendTo != nil {
		sendMap := make(map[string][]string)
		p.service.GetSendMap(body.SendTo, sendMap)

		ch := make(chan error)
		var wg sync.WaitGroup

		for sendUrl, urls := range sendMap {
			wg.Add(1)
			body := types.PeerPresentationBody{
				NewPeer: newPeer,
				SendTo:  urls,
			}
			go p.service.SendNewPeer(body, sendUrl, ch, &wg)
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for err := range ch {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			}
		}
	}

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

	//validate restaurant
	newRestaurant := models.Restaurant{}
	err := c.BodyParser(newRestaurant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	err = p.restaurants.CompleteRestaurantInitialData(&newRestaurant)
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

	// get restaurant ✅
	// complete data ✅
	// validate is in influence area ✅
	// validate is unique (others peers) ⬅
	// save restaurant ✅
	// return restaurant

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}

func (p *PeerController) HaveRestaurant(c *fiber.Ctx) error {
	restaurantQuery := make(map[string]interface{})
	err := c.BodyParser(restaurantQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	have, err := p.service.HaveRestaurant(restaurantQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": have})
}
