package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func CacheObject(key string, obj interface{}) {
	err := redisClient.Set(context.Background(), key, obj, 0).Err()
	if err != nil {
		log.Fatalf("Could not write object to Redis cache, %v", err)
	}
}

func DeleteKey(ctx context.Context, key string) {
	redisClient.Del(ctx, key)
}
