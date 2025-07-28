package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
	Cache      CacheConfig      `mapstructure:"cache"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string
	Password string
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

type BlockchainConfig struct {
	Ethereum ChainConfig `mapstructure:"ethereum"`
	BSC      ChainConfig `mapstructure:"bsc"`
	Polygon  ChainConfig `mapstructure:"polygon"`
}

type ChainConfig struct {
	RPCURL string `mapstructure:"rpc_url"`
}

type CacheConfig struct {
	TokenBalanceTTL string `mapstructure:"token_balance_ttl"`
}

func LoadConfig() (*Config, error) {
	// 加载 .env 文件
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 从环境变量获取敏感信息
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")

	return &config, nil
}
