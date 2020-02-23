package newsrepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	newsmodel "github.com/bungysheep/news-consumer/pkg/models/v1/news"
	elasticv7 "github.com/olivere/elastic/v7"
)

// INewsRepository type
type INewsRepository interface {
	SaveRecord(context.Context, *newsmodel.News) (int64, error)
	SaveNewsID(context.Context, *newsmodel.News) error
}

type newsRepository struct {
	DB       *sql.DB
	ESClient *elasticv7.Client
}

// NewNewsRepository - Create news repository
func NewNewsRepository(db *sql.DB, esClient *elasticv7.Client) INewsRepository {
	return &newsRepository{
		DB:       db,
		ESClient: esClient,
	}
}

// SaveRecord - Save news record into database
func (newsRepo *newsRepository) SaveRecord(ctx context.Context, data *newsmodel.News) (int64, error) {
	conn, err := newsRepo.DB.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO news (author, body, created) 
		VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert news, error: %v", err)
	}
	defer stmt.Close()

	var lastInsertID int64
	err = stmt.QueryRowContext(ctx, data.GetAuthor(), data.GetBody(), data.GetCreated()).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("Failed inserting news, error: %v", err)
	}

	return lastInsertID, nil
}

// SaveNewsID - Save news id and created into elasticsearch
func (newsRepo *newsRepository) SaveNewsID(ctx context.Context, data *newsmodel.News) error {
	dataIndex := map[string]interface{}{
		"id":      data.GetID(),
		"created": data.GetCreated().Format(time.RFC3339),
	}

	dataJSON, err := json.Marshal(dataIndex)
	if err != nil {
		return err
	}

	_, err = newsRepo.ESClient.Index().Index("news").BodyJson(string(dataJSON)).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
