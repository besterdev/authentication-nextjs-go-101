package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
)

const userKeyPrefix = "user:"

type RedisUserCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisUserCache(client *redis.Client, ttl time.Duration) *RedisUserCache {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &RedisUserCache{client: client, ttl: ttl}
}

func (c *RedisUserCache) Get(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	data, err := c.client.Get(ctx, userKeyPrefix+userID.String()).Bytes()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *RedisUserCache) Set(ctx context.Context, user *models.User) error {
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, userKeyPrefix+user.ID.String(), payload, c.ttl).Err()
}

func (c *RedisUserCache) Delete(ctx context.Context, userID uuid.UUID) error {
	err := c.client.Del(ctx, userKeyPrefix+userID.String()).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}
