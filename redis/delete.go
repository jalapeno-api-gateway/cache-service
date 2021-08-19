package redis

import "context"

func DeleteKey(ctx context.Context, key string) {
	redisClient.Del(ctx, key)
}
