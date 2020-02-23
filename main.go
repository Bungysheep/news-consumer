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
	"github.com/bungysheep/news-consumer/pkg/protocols/mq"
	"github.com/bungysheep/news-consumer/pkg/protocols/redis"
	"github.com/bungysheep/news-consumer/pkg/repositories/v1/newsrepository"
	"github.com/bungysheep/news-consumer/pkg/services/v1/newsservice"
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

	if err := mq.CreateMqConnection(); err != nil {
		return err
	}

	if err := elasticsearch.CreateESClient(); err != nil {
		return err
	}

	if err := database.CreateDbConnection(); err != nil {
		return err
	}

	newsSvc := newsservice.NewNewsService(newsrepository.NewNewsRepository(database.DbConnection, elasticsearch.ESClient))

	mqChan, err := mq.MqConnection.Channel()
	if err != nil {
		return err
	}
	defer mqChan.Close()

	queue, err := mqChan.QueueDeclare(configs.MQNEWSPOSTQUEUE, false, false, false, false, nil)
	if err != nil {
		return err
	}

	messages, err := mqChan.Consume(queue.Name, "", true, false, false, false, nil)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for d := range messages {
			log.Printf("Received %v", string(d.Body))
			news := news.NewNews()
			if err := json.Unmarshal(d.Body, &news); err != nil {
				log.Printf("Failed to unmarshal message, error: %v", err)
			}

			if err := newsSvc.DoSave(ctx, news); err != nil {
				log.Printf("Failed to save news, error: %v", err)
			}
		}
	}()

	<-c

	log.Printf("Closing redis client...\n")
	redis.RedisClient.Close()

	log.Printf("Closing rabbitmq connection...\n")
	mq.MqConnection.Close()

	log.Printf("Closing database connection...\n")
	database.DbConnection.Close()

	return nil
}
