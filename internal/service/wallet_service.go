package service

import (
	"wallet-tracker/internal/model"
	"wallet-tracker/internal/repository"
	"wallet-tracker/pkg/cache"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
	cache      *cache.RedisClient
}

func NewWalletService(walletRepo *repository.WalletRepository, cache *cache.RedisClient) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		cache:      cache,
	}
}

func (ws *WalletService) AddWallet(userID uint, address string, chainID int, chainName, name string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID:    userID,
		Address:   address,
		ChainID:   chainID,
		ChainName: chainName,
		Name:      name,
	}

	return ws.walletRepo.Create(wallet)
}

func (ws *WalletService) AddTokenToWallet(walletID uint, tokenAddress, symbol, name string, decimals int) (*model.WalletToken, error) {
	token := &model.WalletToken{
		WalletID:      walletID,
		TokenAddress:  tokenAddress,
		TokenSymbol:   symbol,
		TokenName:     name,
		TokenDecimals: decimals,
		IsActive:      true,
	}

	return ws.walletRepo.CreateToken(token)
}

func (ws *WalletService) GetUserWallets(userID uint) ([]model.Wallet, error) {
	return ws.walletRepo.GetByUserID(userID)
}

func (ws *WalletService) RefreshUserCache(userID uint) error {
	wallets, err := ws.walletRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	var addresses []string
	for _, wallet := range wallets {
		addresses = append(addresses, wallet.Address)
	}

	return ws.cache.DeleteUserBalances(addresses)
}
