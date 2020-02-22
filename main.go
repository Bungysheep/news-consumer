package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bungysheep/news-consumer/pkg/configs"
	"github.com/bungysheep/news-consumer/pkg/protocols/database"
	"github.com/bungysheep/news-consumer/pkg/protocols/elasticsearch"
	"github.com/bungysheep/news-consumer/pkg/protocols/redis"
	_ "github.com/lib/pq"
)

func main() {
	if err := startUp(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func startUp() error {
	if err := redis.CreateRedisClient(); err != nil {
		return err
	}

	if err := elasticsearch.CreateESClient(); err != nil {
		return err
	}

	if err := database.CreateDbConnection(); err != nil {
		return err
	}

	pubSub := redis.RedisClient.Subscribe(configs.REDISNEWSPOSTCHANNEL)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			data, err := pubSub.Receive()
			if err != nil {
				log.Fatalf("Failed to subscribe news channel, error: %v", err)
			}
			log.Printf("%v", data)
		}
	}()

	<-c

	if err := pubSub.Close(); err != nil {
		log.Fatalf("Failed to close subscriber news channel, error: %v", err)
	}

	log.Printf("Closing redis client...\n")
	redis.RedisClient.Close()

	log.Printf("Closing database connection...\n")
	database.DbConnection.Close()

	return nil
}
