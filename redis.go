package samplelib

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

func SetItem(rdb *redis.Client, key string, i *Item) error {
	ctx := context.Background()
	item, err := json.Marshal(i)
	if err != nil {
		return err
	}
	rdb.Set(ctx, key, string(item), 0).Err()
}

func GetItem(rdb *redis.Client, key string) (string, error) {
	ctx := context.Background()
	rdb.Get(ctx, key).Result()
}
