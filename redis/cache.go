package redis

import (
	"context"
	"log"
)

func CacheObject(key string, obj interface{}) {
	err := redisClient.Set(context.Background(), key, obj, 0).Err()
	if err != nil {
		log.Fatalf("Could not write object to Redis cache, %v", err)
	}
}