package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/bcryptpool"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/cache"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/config"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/dto"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/repository"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthService struct {
	cfg           *config.Config
	users         *repository.UserRepository
	refreshTokens *repository.RefreshTokenRepository
	refreshStore  cache.RefreshTokenStore
	userCache     cache.UserCache
	bcryptPool    *bcryptpool.Pool
}

func NewAuthService(
	cfg *config.Config,
	users *repository.UserRepository,
	refreshTokens *repository.RefreshTokenRepository,
	refreshStore cache.RefreshTokenStore,
	userCache cache.UserCache,
	bcryptPool *bcryptpool.Pool,
) *AuthService {
	if refreshStore == nil {
		refreshStore = cache.NewPGRefreshTokenStore(refreshTokens)
	}
	return &AuthService{
		cfg:           cfg,
		users:         users,
		refreshTokens: refreshTokens,
		refreshStore:  refreshStore,
		userCache:     userCache,
		bcryptPool:    bcryptPool,
	}
}

type AccessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	hash, err := s.bcryptPool.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.users.Create(user); err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	if s.userCache != nil {
		_ = s.userCache.Set(context.Background(), user)
	}

	return toUserResponse(user), nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.bcryptPool.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(user)
}

func (s *AuthService) Refresh(refreshToken string) (*dto.AuthResponse, error) {
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	tokenHash := hashToken(refreshToken)
	stored, err := s.refreshStore.FindByTokenHash(context.Background(), tokenHash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if stored.UserID != claims.UserID {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.findUser(claims.UserID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	return s.rotateAndIssueTokens(stored.ID, user)
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.refreshStore.DeleteByUserID(context.Background(), userID)
}

func (s *AuthService) GetMe(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.findUser(userID)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *AuthService) VerifyAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.cfg.JWTAccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthService) findUser(userID uuid.UUID) (*models.User, error) {
	if s.userCache != nil {
		if user, err := s.userCache.Get(context.Background(), userID); err == nil {
			return user, nil
		}
	}

	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if s.userCache != nil {
		_ = s.userCache.Set(context.Background(), user)
	}

	return user, nil
}

func (s *AuthService) issueTokens(user *models.User) (*dto.AuthResponse, error) {
	now := time.Now()
	accessExp := now.Add(s.cfg.JWTAccessExpiry)
	refreshExp := now.Add(s.cfg.JWTRefreshExpiry)

	accessToken, err := s.signAccessToken(user, now, accessExp)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshRecord, err := s.newRefreshToken(user, now, refreshExp)
	if err != nil {
		return nil, err
	}

	if err := s.refreshStore.Create(context.Background(), refreshRecord); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
	}, nil
}

func (s *AuthService) rotateAndIssueTokens(oldID uuid.UUID, user *models.User) (*dto.AuthResponse, error) {
	now := time.Now()
	accessExp := now.Add(s.cfg.JWTAccessExpiry)
	refreshExp := now.Add(s.cfg.JWTRefreshExpiry)

	accessToken, err := s.signAccessToken(user, now, accessExp)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshRecord, err := s.newRefreshToken(user, now, refreshExp)
	if err != nil {
		return nil, err
	}

	if err := s.refreshStore.Rotate(context.Background(), oldID, refreshRecord); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
	}, nil
}

func (s *AuthService) signAccessToken(user *models.User, now, accessExp time.Time) (string, error) {
	accessClaims := AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.cfg.JWTAccessSecret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return accessToken, nil
}

func (s *AuthService) newRefreshToken(user *models.User, now, refreshExp time.Time) (string, *models.RefreshToken, error) {
	refreshClaims := RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.cfg.JWTRefreshSecret))
	if err != nil {
		return "", nil, fmt.Errorf("sign refresh token: %w", err)
	}

	refreshRecord := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: refreshExp,
	}

	return refreshToken, refreshRecord, nil
}

func (s *AuthService) parseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.cfg.JWTRefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
