package middleware

import (
	"net"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	redisstorage "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
)

func AuthRateLimiter(redisClient *redis.Client) fiber.Handler {
	host, portStr, err := net.SplitHostPort(redisClient.Options().Addr)
	if err != nil {
		host = redisClient.Options().Addr
		portStr = "6379"
	}
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 6379
	}

	storage := redisstorage.New(redisstorage.Config{
		Host:     host,
		Port:     port,
		Password: redisClient.Options().Password,
		Database: redisClient.Options().DB,
	})

	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
				"code":  "RATE_LIMITED",
			})
		},
	})
}

func InMemoryAuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
				"code":  "RATE_LIMITED",
			})
		},
	})
}
