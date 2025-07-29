package service

import (
	"errors"
	"os"
	"testing"

	"wallet-tracker/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("Successful Registration", func(t *testing.T) {
		// 模拟用户名和邮箱不存在
		mockRepo.On("GetByUsername", "newuser").Return(nil, errors.New("not found")).Once()
		mockRepo.On("GetByEmail", "new@example.com").Return(nil, errors.New("not found")).Once()

		// 模拟创建成功
		expectedUser := &model.User{
			ID:       1,
			Username: "newuser",
			Email:    "new@example.com",
		}
		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(expectedUser, nil).Once()

		user, err := service.Register("newuser", "new@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "newuser", user.Username)
		assert.Equal(t, "new@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Username Already Exists", func(t *testing.T) {
		existingUser := &model.User{Username: "existinguser"}
		mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil).Once()

		user, err := service.Register("existinguser", "new@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "username already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// 设置 JWT secret
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("Successful Login", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		existingUser := &model.User{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByUsername", "testuser").Return(existingUser, nil).Once()

		token, user, err := service.Login("testuser", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("not found")).Once()

		token, user, err := service.Login("nonexistent", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		existingUser := &model.User{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByUsername", "testuser").Return(existingUser, nil).Once()

		token, user, err := service.Login("testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})
}
