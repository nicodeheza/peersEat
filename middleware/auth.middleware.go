package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/nicodeheza/peersEat/services"
	"github.com/nicodeheza/peersEat/types"
)

type AuthMiddleware struct {
	store             *session.Store
	restaurantService services.RestaurantServiceI
}

type AuthMiddlewareI interface {
	Sessions(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Authenticate(c *fiber.Ctx) error
	Protect(c *fiber.Ctx) error
	OnlyPeerOwner(c *fiber.Ctx) error
}

func InitAuthMiddleware(restaurantService services.RestaurantServiceI) *AuthMiddleware {
	storage := redis.New(redis.Config{
		URL:   os.Getenv("REDIS_URI"),
		Reset: false,
	})
	store := session.New(session.Config{
		Storage:      storage,
		CookieSecure: true,
	})

	return &AuthMiddleware{store, restaurantService}
}

func (a *AuthMiddleware) Sessions(c *fiber.Ctx) error {
	sess, err := a.store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	id := sess.Get("_id")

	c.Locals("restaurantId", fmt.Sprintf("%v", id))
	return c.Next()
}

func (a *AuthMiddleware) Logout(c *fiber.Ctx) error {
	sess, err := a.store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Next()
}

func (a *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	body := new(types.AuthReq)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	isAuth, id, err := a.restaurantService.Authenticate(body.Password, body.UserName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if !isAuth {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	sess, err := a.store.Get(c)

	sess.Set("_id", id)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Next()
}

func (a *AuthMiddleware) Protect(c *fiber.Ctx) error {
	sess, err := a.store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	id := sess.Get("_id")

	w, ok := id.(string)
	if !ok && len(w) < 4 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	return c.Next()
}

func (a *AuthMiddleware) OnlyPeerOwner(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()

	key := os.Getenv("API_KEY")

	if headers["X-API-KEY"] != key {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	return c.Next()
}
