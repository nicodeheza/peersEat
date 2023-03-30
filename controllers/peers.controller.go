package controllers

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/services/validations"
	"github.com/nicodeheza/peersEat/types"
)

type PeerControllerI interface {
	PeerPresentation(c *fiber.Ctx) error
	SendAllPeers(c *fiber.Ctx) error
}

type PeerController struct {
	service  services.PeerServiceI
	validate validations.ValidateI
}

func NewPeerController(service services.PeerServiceI, validate validations.ValidateI) *PeerController {
	return &PeerController{service, validate}
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
