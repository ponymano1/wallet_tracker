package repository

import (
	"testing"

	"wallet-tracker/internal/model"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(&model.User{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = NewUserRepository(db)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	createdUser, err := suite.repo.Create(user)
	suite.NoError(err)
	suite.NotNil(createdUser)
	suite.NotZero(createdUser.ID)
	suite.Equal("testuser", createdUser.Username)
}

func (suite *UserRepositoryTestSuite) TestGetByUsername() {
	// 先创建用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	_, err := suite.repo.Create(user)
	suite.NoError(err)

	// 测试查找
	foundUser, err := suite.repo.GetByUsername("testuser")
	suite.NoError(err)
	suite.NotNil(foundUser)
	suite.Equal("testuser", foundUser.Username)

	// 测试查找不存在的用户
	notFoundUser, err := suite.repo.GetByUsername("nonexistent")
	suite.Error(err)
	suite.Nil(notFoundUser)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	// 先创建用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	_, err := suite.repo.Create(user)
	suite.NoError(err)

	// 测试查找
	foundUser, err := suite.repo.GetByEmail("test@example.com")
	suite.NoError(err)
	suite.NotNil(foundUser)
	suite.Equal("test@example.com", foundUser.Email)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
