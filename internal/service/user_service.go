package service

import (
	"errors"
	"os"
	"time"

	"wallet-tracker/internal/model"
	"wallet-tracker/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	Register(username, email, password string) (*model.User, error)
	Login(username, password string) (string, *model.User, error)
}

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) Register(username, email, password string) (*model.User, error) {
	// 检查用户是否已存在
	if _, err := us.userRepo.GetByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	if _, err := us.userRepo.GetByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return us.userRepo.Create(user)
}

func (us *UserService) Login(username, password string) (string, *model.User, error) {
	user, err := us.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// 生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
