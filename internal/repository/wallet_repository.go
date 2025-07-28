package repository

import (
	"wallet-tracker/internal/model"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (wr *WalletRepository) Create(wallet *model.Wallet) (*model.Wallet, error) {
	if err := wr.db.Create(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wr *WalletRepository) CreateToken(token *model.WalletToken) (*model.WalletToken, error) {
	if err := wr.db.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (wr *WalletRepository) GetByUserID(userID uint) ([]model.Wallet, error) {
	var wallets []model.Wallet
	if err := wr.db.Preload("Tokens").Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (wr *WalletRepository) GetByID(id uint) (*model.Wallet, error) {
	var wallet model.Wallet
	if err := wr.db.Preload("Tokens").First(&wallet, id).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
