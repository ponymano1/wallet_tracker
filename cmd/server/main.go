package main

import (
	"log"

	"wallet-tracker/internal/config"
	"wallet-tracker/internal/handler"
	"wallet-tracker/internal/middleware"
	"wallet-tracker/internal/repository"
	"wallet-tracker/internal/service"
	"wallet-tracker/pkg/cache"
	"wallet-tracker/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// 连接数据库
	db, err := database.NewMySQLConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// 连接 Redis
	redisClient, err := cache.NewRedisClient(&cfg.Redis, &cfg.Cache)
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	// 初始化 repositories
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// 初始化 services
	userService := service.NewUserService(userRepo)
	walletService := service.NewWalletService(walletRepo, redisClient)
	blockchainService, err := service.NewBlockchainService(&cfg.Blockchain, redisClient)
	if err != nil {
		log.Fatal("Failed to initialize blockchain service: ", err)
	}

	// 初始化 handlers
	authHandler := handler.NewAuthHandler(userService)
	walletHandler := handler.NewWalletHandler(walletService, blockchainService)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/wallets", walletHandler.AddWallet)
		protected.GET("/wallets", walletHandler.GetWallets)
		protected.POST("/wallets/:wallet_id/tokens", walletHandler.AddToken)
		protected.GET("/balances", walletHandler.GetBalances)
		protected.POST("/refresh-cache", walletHandler.RefreshCache)
	}

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
