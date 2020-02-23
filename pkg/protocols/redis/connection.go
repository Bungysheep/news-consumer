package redis

import (
	"log"
	"os"

	"github.com/bungysheep/news-consumer/pkg/configs"
	"github.com/go-redis/redis/v7"
)

var (
	// RedisClient - Redis client
	RedisClient *redis.Client
)

// CreateRedisClient - Creates redis client
func CreateRedisClient() error {
	log.Printf("Creating redis client...")

	redisURL, err := resolveRedisURL()
	if err != nil {
		return err
	}

	redisAuth, err := resolveRedisAuth()
	if err != nil {
		return err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisAuth,
		DB:       0,
	})

	RedisClient = client

	_, err = client.Ping().Result()

	return err
}

func resolveRedisURL() (string, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		return redisURL, nil
	}

	return configs.REDISURL, nil
}

func resolveRedisAuth() (string, error) {
	redisAuth := os.Getenv("REDIS_AUTH")
	if redisAuth != "" {
		return redisAuth, nil
	}

	return configs.REDISAUTH, nil
}
