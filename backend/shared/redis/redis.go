package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Connect() {
	addr := getEnv("REDIS_ADDR", "localhost:6380")
	password := getEnv("REDIS_PASSWORD", "123456")

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
		PoolSize: 1000,
		MinIdleConns: 100,
		PoolTimeout: 2 * time.Second,
	})

	pong, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

func PopSubmission(queueName string, timeout time.Duration) (string, bool) {
	result, err := Client.BRPop(Ctx, timeout, queueName).Result()
	if err != nil {
		return "", false
	}
	return result[1], true
}

func PushResult(queueName string, data string) error {
	return Client.LPush(Ctx, queueName, data).Err()
}