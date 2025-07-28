package service

import (
	"fmt"
	"math/big"

	"wallet-tracker/internal/config"
	"wallet-tracker/internal/model"
	"wallet-tracker/pkg/blockchain"
	"wallet-tracker/pkg/cache"
)

type BlockchainService struct {
	clients map[int]*blockchain.BlockchainClient
	cache   *cache.RedisClient
}

func NewBlockchainService(cfg *config.BlockchainConfig, cache *cache.RedisClient) (*BlockchainService, error) {
	clients := make(map[int]*blockchain.BlockchainClient)

	// 以太坊主网 (Chain ID: 1)
	ethClient, err := blockchain.NewBlockchainClient(cfg.Ethereum.RPCURL)
	if err != nil {
		return nil, err
	}
	clients[1] = ethClient

	// BSC 主网 (Chain ID: 56)
	bscClient, err := blockchain.NewBlockchainClient(cfg.BSC.RPCURL)
	if err != nil {
		return nil, err
	}
	clients[56] = bscClient

	// Polygon 主网 (Chain ID: 137)
	polygonClient, err := blockchain.NewBlockchainClient(cfg.Polygon.RPCURL)
	if err != nil {
		return nil, err
	}
	clients[137] = polygonClient

	return &BlockchainService{
		clients: clients,
		cache:   cache,
	}, nil
}

func (bs *BlockchainService) GetTokenBalance(chainID int, tokenAddress, walletAddress string, forceRefresh bool) (*model.TokenBalance, error) {
	// 如果不强制刷新，先尝试从缓存获取
	if !forceRefresh {
		if cached, err := bs.cache.GetTokenBalance(walletAddress, tokenAddress); err == nil {
			return cached, nil
		}
	}

	client, exists := bs.clients[chainID]
	if !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	// 获取余额
	balance, err := client.GetTokenBalance(tokenAddress, walletAddress)
	if err != nil {
		return nil, err
	}

	// 获取 token 信息
	symbol, name, decimals, err := client.GetTokenInfo(tokenAddress)
	if err != nil {
		return nil, err
	}

	// 计算实际余额（考虑小数位）
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	balanceFloat := new(big.Float).SetInt(balance)
	divisorFloat := new(big.Float).SetInt(divisor)
	actualBalance := new(big.Float).Quo(balanceFloat, divisorFloat)

	chainName := getChainName(chainID)

	tokenBalance := &model.TokenBalance{
		WalletAddress: walletAddress,
		TokenAddress:  tokenAddress,
		Balance:       actualBalance.String(),
		Symbol:        symbol,
		Name:          name,
		Decimals:      int(decimals),
		ChainID:       chainID,
		ChainName:     chainName,
	}

	// 缓存结果
	bs.cache.SetTokenBalance(walletAddress, tokenAddress, tokenBalance)

	return tokenBalance, nil
}

func (bs *BlockchainService) GetMultipleTokenBalances(wallets []model.Wallet, forceRefresh bool) ([]model.TokenBalance, error) {
	var results []model.TokenBalance

	for _, wallet := range wallets {
		for _, token := range wallet.Tokens {
			if !token.IsActive {
				continue
			}

			balance, err := bs.GetTokenBalance(wallet.ChainID, token.TokenAddress, wallet.Address, forceRefresh)
			if err != nil {
				continue // 跳过错误的 token，继续处理其他的
			}

			results = append(results, *balance)
		}
	}

	return results, nil
}

func getChainName(chainID int) string {
	switch chainID {
	case 1:
		return "Ethereum"
	case 56:
		return "BSC"
	case 137:
		return "Polygon"
	default:
		return "Unknown"
	}
}
