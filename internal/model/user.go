package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Wallets   []Wallet       `json:"wallets" gorm:"foreignKey:UserID"`
}

type Wallet struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Address   string         `json:"address" gorm:"not null"`
	ChainID   int            `json:"chain_id" gorm:"not null"`
	ChainName string         `json:"chain_name" gorm:"not null"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Tokens    []WalletToken  `json:"tokens" gorm:"foreignKey:WalletID"`
}

type WalletToken struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	WalletID      uint           `json:"wallet_id" gorm:"not null"`
	TokenAddress  string         `json:"token_address" gorm:"not null"`
	TokenSymbol   string         `json:"token_symbol"`
	TokenName     string         `json:"token_name"`
	TokenDecimals int            `json:"token_decimals"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type TokenBalance struct {
	WalletAddress string  `json:"wallet_address"`
	TokenAddress  string  `json:"token_address"`
	Balance       string  `json:"balance"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Decimals      int     `json:"decimals"`
	ChainID       int     `json:"chain_id"`
	ChainName     string  `json:"chain_name"`
	USDValue      float64 `json:"usd_value,omitempty"`
}
