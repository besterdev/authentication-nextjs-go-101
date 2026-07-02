package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/repository"
)

const refreshKeyPrefix = "refresh:"

var rotateScript = redis.NewScript(`
local key = KEYS[1]
if redis.call("EXISTS", key) == 0 then
  return 0
end
redis.call("DEL", key)
redis.call("SET", ARGV[1], ARGV[2], "EX", ARGV[3])
return 1
`)

type RedisRefreshTokenStore struct {
	client *redis.Client
	repo   *repository.RefreshTokenRepository
}

func NewRedisRefreshTokenStore(client *redis.Client, repo *repository.RefreshTokenRepository) *RedisRefreshTokenStore {
	return &RedisRefreshTokenStore{client: client, repo: repo}
}

func (s *RedisRefreshTokenStore) Create(ctx context.Context, token *models.RefreshToken) error {
	if err := s.repo.Create(token); err != nil {
		return err
	}
	return s.setRedis(ctx, token)
}

func (s *RedisRefreshTokenStore) FindByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	data, err := s.client.Get(ctx, refreshKeyPrefix+tokenHash).Bytes()
	if err == nil {
		var token models.RefreshToken
		if json.Unmarshal(data, &token) == nil {
			return &token, nil
		}
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	token, err := s.repo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}

	_ = s.setRedis(ctx, token)
	return token, nil
}

func (s *RedisRefreshTokenStore) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *RedisRefreshTokenStore) Rotate(ctx context.Context, oldID uuid.UUID, newToken *models.RefreshToken) error {
	oldToken, err := s.repo.FindByID(oldID)
	if err != nil {
		return err
	}

	if err := s.repo.Rotate(oldID, newToken); err != nil {
		return err
	}

	oldKey := refreshKeyPrefix + oldToken.TokenHash
	newKey := refreshKeyPrefix + newToken.TokenHash
	payload, err := json.Marshal(newToken)
	if err != nil {
		return err
	}

	ttl := time.Until(newToken.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("refresh token already expired")
	}

	_, err = rotateScript.Run(ctx, s.client, []string{oldKey}, newKey, payload, int(ttl.Seconds())).Int()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}

func (s *RedisRefreshTokenStore) setRedis(ctx context.Context, token *models.RefreshToken) error {
	payload, err := json.Marshal(token)
	if err != nil {
		return err
	}
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return nil
	}
	return s.client.Set(ctx, refreshKeyPrefix+token.TokenHash, payload, ttl).Err()
}
