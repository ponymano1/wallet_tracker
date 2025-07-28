package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"wallet-tracker/internal/config"
	"wallet-tracker/internal/model"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

func NewRedisClient(cfg *config.RedisConfig, cacheCfg *config.CacheConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(cacheCfg.TokenBalanceTTL)
	if err != nil {
		ttl = 24 * time.Hour // 默认24小时
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
		ttl:    ttl,
	}, nil
}

func (r *RedisClient) SetTokenBalance(walletAddress, tokenAddress string, balance *model.TokenBalance) error {
	key := fmt.Sprintf("balance:%s:%s", walletAddress, tokenAddress)
	data, err := json.Marshal(balance)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, r.ttl).Err()
}

func (r *RedisClient) GetTokenBalance(walletAddress, tokenAddress string) (*model.TokenBalance, error) {
	key := fmt.Sprintf("balance:%s:%s", walletAddress, tokenAddress)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var balance model.TokenBalance
	if err := json.Unmarshal([]byte(data), &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

func (r *RedisClient) DeleteTokenBalance(walletAddress, tokenAddress string) error {
	key := fmt.Sprintf("balance:%s:%s", walletAddress, tokenAddress)
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) DeleteUserBalances(walletAddresses []string) error {
	var keys []string
	for _, address := range walletAddresses {
		pattern := fmt.Sprintf("balance:%s:*", address)
		matchedKeys, err := r.client.Keys(r.ctx, pattern).Result()
		if err != nil {
			continue
		}
		keys = append(keys, matchedKeys...)
	}

	if len(keys) > 0 {
		return r.client.Del(r.ctx, keys...).Err()
	}
	return nil
}
