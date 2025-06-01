package redisutil

import (
	"context"
	"encoding/json"

	"github.com/typical-developers/discord-bot-backend/internal/cache"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func GetCached[T any](ctx context.Context, key string) *T {
	v := cache.Client.JSONGet(ctx, key, "$")

	if v.Val() == "" {
		return nil
	}

	expanded, err := v.Expanded()
	if err != nil {
		logger.Log.Error("Failed to expand cached value: %s", "error", err)
		return nil
	}

	entries := expanded.([]any)
	entry := entries[0].(map[string]any)

	if entry == nil {
		return nil
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		logger.Log.Error("Failed to marshal cached data: %s", "error", err)
		return nil
	}

	var t T
	err = json.Unmarshal(jsonBytes, &t)
	if err != nil {
		logger.Log.Error("Failed to unmarshal cached data: %s", "error", err)
		return nil
	}

	return &t
}

func SetCached[T any](ctx context.Context, key string, value T, opts *CacheOpts) {
	cache.Client.JSONSet(ctx, key, "$", value)

	if opts.Expiry != nil {
		_ = cache.Client.Expire(ctx, key, *opts.Expiry)
	}
}
