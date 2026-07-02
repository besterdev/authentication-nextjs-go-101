package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
	"gorm.io/gorm"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepository) FindByID(id uuid.UUID) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("id = ?", id).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) FindByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) DeleteByID(id uuid.UUID) error {
	return r.db.Delete(&models.RefreshToken{}, "id = ?", id).Error
}

func (r *RefreshTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&models.RefreshToken{}, "user_id = ?", userID).Error
}

func (r *RefreshTokenRepository) DeleteExpired(before time.Time) (int64, error) {
	result := r.db.Where("expires_at <= ?", before).Delete(&models.RefreshToken{})
	return result.RowsAffected, result.Error
}

func (r *RefreshTokenRepository) Rotate(oldID uuid.UUID, newToken *models.RefreshToken) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.RefreshToken{}, "id = ?", oldID).Error; err != nil {
			return err
		}
		return tx.Create(newToken).Error
	})
}
