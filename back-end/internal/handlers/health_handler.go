package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/cache"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redisClient}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx := c.Context()

	if err := h.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"db":     "down",
		})
	}

	redisStatus := "disabled"
	if h.redis != nil {
		if err := cache.Ping(ctx, h.redis); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"db":     "up",
				"redis":  "down",
			})
		}
		redisStatus = "up"
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"db":     "up",
		"redis":  redisStatus,
	})
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return h.Check(c)
}

func PingDB(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Exec("SELECT 1").Error
}
