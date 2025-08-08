package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"ua/shared/logger"
	"go.uber.org/zap"
)

type Client struct {
	*redis.Client
}

func NewRedisClient(redisURL string) (*Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Info("Connected to Redis")
	return &Client{rdb}, nil
}

func (c *Client) Close() error {
	logger.Info("Closing Redis connection")
	return c.Client.Close()
}

func (c *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	logger.Debug("Setting JSON value in Redis", zap.String("key", key))
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) GetJSON(ctx context.Context, key string) (string, error) {
	logger.Debug("Getting JSON value from Redis", zap.String("key", key))
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	logger.Debug("Adding to sorted set", zap.String("key", key))
	return c.Client.ZAdd(ctx, key, members...).Err()
}

func (c *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	logger.Debug("Getting range by score", zap.String("key", key))
	return c.Client.ZRangeByScore(ctx, key, opt).Result()
}

func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	logger.Debug("Pushing to list", zap.String("key", key))
	return c.Client.LPush(ctx, key, values...).Err()
}

func (c *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	logger.Debug("Getting list range", zap.String("key", key))
	return c.Client.LRange(ctx, key, start, stop).Result()
}