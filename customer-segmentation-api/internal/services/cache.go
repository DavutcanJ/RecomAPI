package services

import (
    "context"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
    addr := os.Getenv("REDIS_URI")
    if addr == "" {
        addr = "localhost:6379"
    }
    
    password := os.Getenv("REDIS_PASSWORD")
    
    db := 0
    if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
        if parsed, err := strconv.Atoi(dbStr); err == nil {
            db = parsed
        }
    }
    
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := RedisClient.Ping(ctx).Err(); err != nil {
        log.Printf("Failed to connect to Redis: %v", err)
        return err
    }
    
    log.Println("Connected to Redis successfully!")
    return nil
}

func GetSegmentFromCache(ctx context.Context, key string) (string, error) {
    if RedisClient == nil {
        return "", redis.Nil // Cache not available
    }
    return RedisClient.Get(ctx, key).Result()
}

func SetSegmentToCache(ctx context.Context, key, value string, ttl time.Duration) error {
    if RedisClient == nil {
        return nil // Cache not available, no error
    }
    return RedisClient.Set(ctx, key, value, ttl).Err()
}