package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

func NewUserRepository(writeDB, readDB *gorm.DB) *UserRepository {
	if readDB == nil {
		readDB = writeDB
	}
	return &UserRepository{writeDB: writeDB, readDB: readDB}
}

func (r *UserRepository) Create(user *models.User) error {
	err := r.writeDB.Create(user).Error
	if isUniqueViolation(err) {
		return ErrEmailAlreadyExists
	}
	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.readDB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.readDB.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
