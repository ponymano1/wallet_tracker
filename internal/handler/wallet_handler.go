package handler

import (
	"net/http"
	"strconv"

	"wallet-tracker/internal/service"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService     *service.WalletService
	blockchainService *service.BlockchainService
}

func NewWalletHandler(walletService *service.WalletService, blockchainService *service.BlockchainService) *WalletHandler {
	return &WalletHandler{
		walletService:     walletService,
		blockchainService: blockchainService,
	}
}

func (wh *WalletHandler) AddWallet(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Address   string `json:"address" binding:"required"`
		ChainID   int    `json:"chain_id" binding:"required"`
		ChainName string `json:"chain_name" binding:"required"`
		Name      string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := wh.walletService.AddWallet(userID, req.Address, req.ChainID, req.ChainName, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

func (wh *WalletHandler) AddToken(c *gin.Context) {
	walletIDStr := c.Param("wallet_id")
	walletID, err := strconv.ParseUint(walletIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	var req struct {
		TokenAddress string `json:"token_address" binding:"required"`
		Symbol       string `json:"symbol" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Decimals     int    `json:"decimals" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := wh.walletService.AddTokenToWallet(uint(walletID), req.TokenAddress, req.Symbol, req.Name, req.Decimals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (wh *WalletHandler) GetWallets(c *gin.Context) {
	userID := c.GetUint("user_id")

	wallets, err := wh.walletService.GetUserWallets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

func (wh *WalletHandler) GetBalances(c *gin.Context) {
	userID := c.GetUint("user_id")
	forceRefresh := c.Query("force_refresh") == "true"

	wallets, err := wh.walletService.GetUserWallets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balances, err := wh.blockchainService.GetMultipleTokenBalances(wallets, forceRefresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balances": balances,
		"cached":   !forceRefresh,
	})
}

func (wh *WalletHandler) RefreshCache(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := wh.walletService.RefreshUserCache(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cache refreshed successfully"})
}
