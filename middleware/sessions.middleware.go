package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

type SessionMiddleware struct {
	store *session.Store
}

func InitSessionMiddleware() *SessionMiddleware {
	storage := redis.New(redis.Config{
		URL:   os.Getenv("REDIS_URI"),
		Reset: false,
	})
	store := session.New(session.Config{
		Storage:      storage,
		CookieSecure: true,
	})

	return &SessionMiddleware{store}
}

func (s *SessionMiddleware) Middleware(c *fiber.Ctx) error {
	sess, err := s.store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	id := sess.Get("_id")

	c.Set("restaurantId", fmt.Sprintf("%v", id))
	return c.Next()
}
