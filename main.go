package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Use environment variables or defaults
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	pass := os.Getenv("REDIS_PASS") // leave empty if none

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Redis connection failed:", err)
		return
	}
	fmt.Println("Redis connected successfully!")

	// Set a key
	if err := rdb.Set(ctx, "testkey", "hello-redis", 0).Err(); err != nil {
		fmt.Println("Failed to set key:", err)
		return
	}
	fmt.Println("Key set successfully!")

	// Get the key
	val, err := rdb.Get(ctx, "testkey").Result()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}
	fmt.Println("Key value:", val)
}
