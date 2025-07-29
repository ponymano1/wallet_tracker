package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-tracker/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(username, email, password string) (*model.User, error) {
	args := m.Called(username, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(username, password string) (string, *model.User, error) {
	args := m.Called(username, password)
	return args.String(0), args.Get(1).(*model.User), args.Error(2)
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Register(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)
	router := setupGin()

	router.POST("/register", handler.Register)

	t.Run("Successful Registration", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockService.On("Register", "testuser", "test@example.com", "password123").
			Return(expectedUser, nil).Once()

		requestBody := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		requestBody := map[string]string{
			"username": "testuser",
			// missing email and password
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)
	router := setupGin()

	router.POST("/login", handler.Login)

	t.Run("Successful Login", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}
		expectedToken := "jwt.token.here"

		mockService.On("Login", "testuser", "password123").
			Return(expectedToken, expectedUser, nil).Once()

		requestBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, response["token"])

		mockService.AssertExpectations(t)
	})
}
