package mock

import (
	"errors"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
)

const errMsg = "example error msg"

type Storage struct {
	GetFetchersErr   bool
	AddFetcherErr    bool
	DeleteFetcherErr bool
	GetHistoryErr    bool
}

func (s Storage) GetFetchers() ([]models.Fetcher, error) {
	if s.GetFetchersErr {
		return nil, errors.New(errMsg)
	}
	return []models.Fetcher{}, nil
}

func (s Storage) AddFetcher(fetcher *models.Fetcher) error {
	if s.AddFetcherErr {
		return errors.New(errMsg)
	}
	return nil
}

func (s Storage) DeleteFetcher(id int) error {
	if s.DeleteFetcherErr {
		return errors.New(errMsg)
	}
	return nil
}

func (s Storage) GetHistory(id int) ([]models.History, error) {
	if s.GetHistoryErr {
		return nil, errors.New(errMsg)
	}
	return []models.History{}, nil
}
