package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&User{}, &Wallet{}, &WalletToken{})
	assert.NoError(t, err)

	return db
}

func TestUserModel(t *testing.T) {
	fmt.Println("TestUserModel")
	db := setupTestDB(t)

	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	// 测试创建用户
	err := db.Create(user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)

	// 测试查找用户
	var foundUser User
	err = db.Where("username = ?", "testuser").First(&foundUser).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestWalletModel(t *testing.T) {
	db := setupTestDB(t)

	// 创建用户
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	// 创建钱包
	wallet := &Wallet{
		UserID:    user.ID,
		Address:   "0x742d35Cc6634C0532925a3b8D2d291b8F0932C71",
		ChainID:   1,
		ChainName: "Ethereum",
		Name:      "Test Wallet",
	}

	err = db.Create(wallet).Error
	assert.NoError(t, err)
	assert.NotZero(t, wallet.ID)

	// 测试关联查询
	var userWithWallets User
	err = db.Preload("Wallets").First(&userWithWallets, user.ID).Error
	assert.NoError(t, err)
	assert.Len(t, userWithWallets.Wallets, 1)
	assert.Equal(t, wallet.Address, userWithWallets.Wallets[0].Address)
}
