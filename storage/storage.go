package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BarTar213/bartlomiej-tarczynski/config"
	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/go-pg/pg/v10"
)

type Storage interface {
	GetFetchers() ([]models.Fetcher, error)
	GetFetcherJob(id int) (int, error)
	AddFetcher(fetcher *models.Fetcher) error
	UpdateFetcher(fetcher *models.Fetcher) error
	UpdateFetcherJobId(fetcherId, jobId int) error
	DeleteFetcher(id int) (int, error)

	GetHistory(id int) ([]models.History, error)
	AddHistory(history *models.History) error
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

	err = initTables(db)
	if err != nil {
		return nil, fmt.Errorf("could not init tables, err: %s", err)
	}

	return &Postgres{db: db}, nil
}

func initTables(db *pg.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		b, err := ioutil.ReadFile("tables.sql")
		if err != nil {
			return err
		}

		queries := strings.Split(string(b), ";")
		for _, query := range queries {
			if len(query) == 0 {
				continue
			}
			_, err := tx.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
