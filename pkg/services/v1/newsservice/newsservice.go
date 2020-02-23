package newsservice

import (
	"context"

	newsmodel "github.com/bungysheep/news-consumer/pkg/models/v1/news"
	"github.com/bungysheep/news-consumer/pkg/repositories/v1/newsrepository"
)

// INewsService type
type INewsService interface {
	DoSave(context.Context, *newsmodel.News) error
}

type newsService struct {
	NewsRepository newsrepository.INewsRepository
}

// NewNewsService - Create news service
func NewNewsService(newsRepo newsrepository.INewsRepository) INewsService {
	return &newsService{
		NewsRepository: newsRepo,
	}
}

// DoPost - Post news
func (newsSvc *newsService) DoSave(ctx context.Context, data *newsmodel.News) error {
	lastInsertID, err := newsSvc.NewsRepository.SaveRecord(ctx, data)
	if err != nil {
		return err
	}

	data.ID = lastInsertID

	return newsSvc.NewsRepository.SaveNewsID(ctx, data)
}
