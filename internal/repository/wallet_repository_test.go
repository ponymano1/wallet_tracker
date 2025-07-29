package repository

import (
	"testing"

	"wallet-tracker/internal/model"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type WalletRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	repo     *WalletRepository
	userRepo *UserRepository
	testUser *model.User
}

func (suite *WalletRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.WalletToken{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = NewWalletRepository(db)
	suite.userRepo = NewUserRepository(db)

	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	createdUser, err := suite.userRepo.Create(user)
	suite.Require().NoError(err)
	suite.testUser = createdUser
}

func (suite *WalletRepositoryTestSuite) TestCreateWallet() {
	wallet := &model.Wallet{
		UserID:    suite.testUser.ID,
		Address:   "0x742d35Cc6634C0532925a3b8D2d291b8F0932C71",
		ChainID:   1,
		ChainName: "Ethereum",
		Name:      "Test Wallet",
	}

	createdWallet, err := suite.repo.Create(wallet)
	suite.NoError(err)
	suite.NotNil(createdWallet)
	suite.NotZero(createdWallet.ID)
	suite.Equal(suite.testUser.ID, createdWallet.UserID)
}

func (suite *WalletRepositoryTestSuite) TestGetByUserID() {
	// 创建多个钱包
	wallet1 := &model.Wallet{
		UserID:    suite.testUser.ID,
		Address:   "0x742d35Cc6634C0532925a3b8D2d291b8F0932C71",
		ChainID:   1,
		ChainName: "Ethereum",
		Name:      "Ethereum Wallet",
	}

	wallet2 := &model.Wallet{
		UserID:    suite.testUser.ID,
		Address:   "0x851d35Cc6634C0532925a3b8D2d291b8F0932C72",
		ChainID:   56,
		ChainName: "BSC",
		Name:      "BSC Wallet",
	}

	_, err := suite.repo.Create(wallet1)
	suite.NoError(err)
	_, err = suite.repo.Create(wallet2)
	suite.NoError(err)

	// 测试获取用户钱包
	wallets, err := suite.repo.GetByUserID(suite.testUser.ID)
	suite.NoError(err)
	suite.Len(wallets, 2)
}

func (suite *WalletRepositoryTestSuite) TestCreateToken() {
	// 先创建钱包
	wallet := &model.Wallet{
		UserID:    suite.testUser.ID,
		Address:   "0x742d35Cc6634C0532925a3b8D2d291b8F0932C71",
		ChainID:   1,
		ChainName: "Ethereum",
		Name:      "Test Wallet",
	}
	createdWallet, err := suite.repo.Create(wallet)
	suite.NoError(err)

	// 创建代币
	token := &model.WalletToken{
		WalletID:      createdWallet.ID,
		TokenAddress:  "0xA0b86a33E6417C9f3b6C37Bb8E0A8b8BBf9f2f71",
		TokenSymbol:   "USDC",
		TokenName:     "USD Coin",
		TokenDecimals: 6,
		IsActive:      true,
	}

	createdToken, err := suite.repo.CreateToken(token)
	suite.NoError(err)
	suite.NotNil(createdToken)
	suite.NotZero(createdToken.ID)
	suite.Equal("USDC", createdToken.TokenSymbol)
}

func TestWalletRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(WalletRepositoryTestSuite))
}
