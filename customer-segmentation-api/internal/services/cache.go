package services

import (
    "context"

    "time"

    "github.com/redis/go-redis/v9"
)

var RedisClient = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",  // Docker'da container name kullan
})

func GetSegmentFromCache(ctx context.Context, key string) (string, error) {
    return RedisClient.Get(ctx, key).Result()
}

func SetSegmentToCache(ctx context.Context, key, value string, ttl time.Duration) error {
    return RedisClient.Set(ctx, key, value, ttl).Err()
}