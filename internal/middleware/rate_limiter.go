package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        50,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
}