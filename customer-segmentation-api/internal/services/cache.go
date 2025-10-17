package services

import (
    "context"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/redis/go-redis/v9"
)

type CacheConfig struct {
    Addr     string
    Password string
    DB       int
}

func getRedisConfig() *CacheConfig {
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
    
    return &CacheConfig{
        Addr:     addr,
        Password: password,
        DB:       db,
    }
}

var RedisClient *redis.Client

func InitRedis() error {
    config := getRedisConfig()
    
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     config.Addr,
        Password: config.Password,
        DB:       config.DB,
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
    return RedisClient.Get(ctx, key).Result()
}

func SetSegmentToCache(ctx context.Context, key, value string, ttl time.Duration) error {
    return RedisClient.Set(ctx, key, value, ttl).Err()
}