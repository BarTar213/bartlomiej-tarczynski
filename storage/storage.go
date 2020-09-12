package storage

import (
	"context"

	"github.com/BarTar213/bartlomiej-tarczynski/config"
	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/go-pg/pg/v10"
)

type Storage interface {
	GetFetchers() ([]models.Fetcher, error)
	AddFetcher(fetcher *models.Fetcher) error
	DeleteFetcher(id int) error

	GetHistory(id int) ([]models.History, error)
}

type Postgres struct {
	db *pg.DB
}

func NewPostgres(config *config.Postgres) (Storage, error) {
	db := pg.Connect(&pg.Options{
		Addr:     config.Address,
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
	})

	err := db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}
