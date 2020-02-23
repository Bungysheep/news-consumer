package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/bungysheep/news-consumer/pkg/configs"
	"github.com/bungysheep/news-consumer/pkg/models/v1/news"
	"github.com/bungysheep/news-consumer/pkg/protocols/database"
	"github.com/bungysheep/news-consumer/pkg/protocols/elasticsearch"
	"github.com/bungysheep/news-consumer/pkg/protocols/redis"
	"github.com/bungysheep/news-consumer/pkg/repositories/v1/newsrepository"
	"github.com/bungysheep/news-consumer/pkg/services/v1/newsservice"
	redisv7 "github.com/go-redis/redis/v7"
	_ "github.com/lib/pq"
)

func main() {
	if err := startUp(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func startUp() error {
	ctx := context.TODO()

	if err := redis.CreateRedisClient(); err != nil {
		return err
	}

	if err := elasticsearch.CreateESClient(); err != nil {
		return err
	}

	if err := database.CreateDbConnection(); err != nil {
		return err
	}

	newsSvc := newsservice.NewNewsService(newsrepository.NewNewsRepository(database.DbConnection, elasticsearch.ESClient))
	pubSub := redis.RedisClient.Subscribe(configs.REDISNEWSPOSTCHANNEL)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			msgi, err := pubSub.Receive()
			if err != nil {
				log.Fatalf("Failed to subscribe news channel, error: %v", err)
			}

			switch msg := msgi.(type) {
			case *redisv7.Subscription:
				log.Printf("Subscribed to %v", msg.Channel)

			case *redisv7.Message:
				log.Printf("Received %v from %v", msg.Payload, msg.Channel)

				if msg.Channel == configs.REDISNEWSPOSTCHANNEL {
					news := news.NewNews()
					if err := json.Unmarshal([]byte(msg.Payload), &news); err != nil {
						log.Printf("Failed to unmarshal payload, error: %v", err)
					}

					if err := newsSvc.DoSave(ctx, news); err != nil {
						log.Printf("Failed to save news channel, error: %v", err)
					}
				}
			}
		}
	}()

	<-c

	if err := pubSub.Close(); err != nil {
		return err
	}

	log.Printf("Closing redis client...\n")
	redis.RedisClient.Close()

	log.Printf("Closing database connection...\n")
	database.DbConnection.Close()

	return nil
}
