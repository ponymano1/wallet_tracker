package repository

import (
	"wallet-tracker/internal/model"

	"gorm.io/gorm"
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user *model.User) (*model.User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := ur.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
