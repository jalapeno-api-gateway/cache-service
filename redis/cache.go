package redis

import (
	"context"
	"log"
)

func CacheLsNode(key string, document LsNodeDocument) {
	err := redisClient.Set(context.Background(), key, document, 0).Err()
	if err != nil {
		log.Fatalf("Could not write LsNode to Redis Cache, %v", err)
	}
}

func CacheLsLink(key string, document LsLinkDocument) {
	err := redisClient.Set(context.Background(), key, document, 0).Err()
	if err != nil {
		log.Fatalf("Could not write LsLink to Redis Cache, %v", err)
	}
}
