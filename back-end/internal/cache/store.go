package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
)

type RefreshTokenStore interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	Rotate(ctx context.Context, oldID uuid.UUID, newToken *models.RefreshToken) error
}

type UserCache interface {
	Get(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Set(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
