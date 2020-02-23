package newsrepository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	newsmodel "github.com/bungysheep/news-consumer/pkg/models/v1/news"
	"github.com/bungysheep/news-consumer/pkg/protocols/elasticsearch"
)

var (
	ctx  context.Context
	repo INewsRepository
	db   *sql.DB
	mock sqlmock.Sqlmock
	data []*newsmodel.News
)

func TestMain(m *testing.M) {
	ctx = context.TODO()

	db, mock, _ = sqlmock.New()
	defer db.Close()

	repo = NewNewsRepository(db, elasticsearch.ESClient)

	data = append(data, &newsmodel.News{
		ID:      1,
		Author:  "Author A",
		Body:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent sapien mauris, viverra vel egestas sit amet, mattis sed libero. Etiam sed diam et felis venenatis porta. Aliquam semper sem eget lectus tristique vulputate. Aliquam euismod nisi at justo congue tempus. Sed faucibus non sapien sit amet condimentum. Sed rutrum ligula odio, sit amet bibendum diam sagittis a. Phasellus sit amet risus tellus. Ut elementum venenatis arcu vitae vulputate. Nulla venenatis, magna et luctus gravida, mi lorem molestie ipsum, sed malesuada erat justo id nunc. Integer sodales sem ac ipsum dapibus lobortis. Vivamus auctor felis non magna ultricies, laoreet posuere tellus posuere.",
		Created: time.Now(),
	})

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestNewsRepository(t *testing.T) {
	t.Run("Save news", saveNews(ctx))
}

func saveNews(ctx context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Save fail", saveFail(ctx, data[0]))

		t.Run("Save success", saveSuccess(ctx, data[0]))
	}
}

func saveFail(ctx context.Context, input *newsmodel.News) func(t *testing.T) {
	return func(t *testing.T) {
		expQuery := mock.ExpectPrepare(`INSERT INTO news`).ExpectQuery()
		expQuery.WithArgs(input.GetAuthor(), input.GetBody(), input.GetCreated()).WillReturnError(fmt.Errorf("Failed saving news"))

		newsID, err := repo.SaveRecord(ctx, input)
		if err == nil {
			t.Errorf("Expect error is not nil")
		}

		if newsID != 0 {
			t.Errorf("Expect news id is 0")
		}
	}
}

func saveSuccess(ctx context.Context, input *newsmodel.News) func(t *testing.T) {
	return func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		expQuery := mock.ExpectPrepare(`INSERT INTO news`).ExpectQuery()
		expQuery.WithArgs(input.GetAuthor(), input.GetBody(), input.GetCreated()).WillReturnRows(rows)

		newsID, err := repo.SaveRecord(ctx, input)
		if err != nil {
			t.Fatalf("Expect error is nil, but got %v", err)
		}

		if newsID != 1 {
			t.Errorf("Expect news id is 1, but got %d", newsID)
		}
	}
}
