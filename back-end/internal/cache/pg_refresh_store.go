package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/repository"
)

type PGRefreshTokenStore struct {
	repo *repository.RefreshTokenRepository
}

func NewPGRefreshTokenStore(repo *repository.RefreshTokenRepository) *PGRefreshTokenStore {
	return &PGRefreshTokenStore{repo: repo}
}

func (s *PGRefreshTokenStore) Create(ctx context.Context, token *models.RefreshToken) error {
	return s.repo.Create(token)
}

func (s *PGRefreshTokenStore) FindByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	return s.repo.FindByTokenHash(tokenHash)
}

func (s *PGRefreshTokenStore) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *PGRefreshTokenStore) Rotate(ctx context.Context, oldID uuid.UUID, newToken *models.RefreshToken) error {
	return s.repo.Rotate(oldID, newToken)
}
