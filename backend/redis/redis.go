package redisClient

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectRedis(addr, username, password string) {
	if addr == "" {
		log.Println("Redis address empty. Continuing without Redis cache.")
		RDB = nil
		return
	}

	opts := &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       0,
	}

	// Use TLS if cloud Redis URL or port 6379/6380
	if username != "" || password != "" {
		opts.TLSConfig = &tls.Config{}
	}

	client := redis.NewClient(opts)

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v (continuing without Redis cache)", err)
		RDB = nil
		return
	}

	RDB = client
	log.Println("Redis connected successfully!")
}