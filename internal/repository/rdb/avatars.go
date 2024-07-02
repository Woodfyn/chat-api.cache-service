package rdb

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Woodfyn/chat-api-cache-service/pkg/core"
	"github.com/redis/go-redis/v9"
)

func (c *Cache) CreateAvatars(ctx context.Context, key string, input []*core.Avatar, ttl time.Duration) error {
	data, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) GetAvatars(ctx context.Context, key string) ([]*core.Avatar, error) {
	bytes, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, core.ErrCacheIsExpiredOrNotFound
		}
		return nil, err
	}

	var result []*core.Avatar
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Cache) UpdateAvatars(ctx context.Context, key string, input []*core.Avatar, ttl time.Duration) error {
	if err := c.rdb.Get(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return core.ErrCacheIsExpiredOrNotFound
		}
		return err
	}

	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}

	data, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, ttl).Err()
}
