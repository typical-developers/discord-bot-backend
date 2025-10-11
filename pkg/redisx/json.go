package redisx

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func JSONUnwrap[T any](ctx context.Context, r *redis.Client, key string, selector string, data *T) error {
	cmd := r.JSONGet(ctx, key, selector)

	err := cmd.Err()
	if err != nil {
		return err
	}

	if cmd.Val() == "" {
		return nil
	}

	expanded, err := cmd.Expanded()
	if err != nil {
		return err
	}

	entries := expanded.([]any)
	entry := entries[0]

	if entry == nil {
		return nil
	}

	jsonB, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonB, &data)
}
