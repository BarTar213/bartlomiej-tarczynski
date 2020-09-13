package mock

import (
	"errors"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
)

const errMsg = "example error msg"

type Storage struct {
	GetFetchersErr        bool
	AddFetcherErr         bool
	UpdateFetcherErr      bool
	UpdateFetcherJobIdErr bool
	DeleteFetcherErr      bool
	GetHistoryErr         bool
	AddHistoryErr         bool
}

func (s *Storage) GetFetchers() ([]models.Fetcher, error) {
	if s.GetFetchersErr {
		return nil, errors.New(errMsg)
	}
	return []models.Fetcher{}, nil
}

func (s *Storage) AddFetcher(fetcher *models.Fetcher) error {
	if s.AddFetcherErr {
		return errors.New(errMsg)
	}
	return nil
}

func (s *Storage) DeleteFetcher(id int) error {
	if s.DeleteFetcherErr {
		return errors.New(errMsg)
	}
	return nil
}

func (s *Storage) GetHistory(id int) ([]models.History, error) {
	if s.GetHistoryErr {
		return nil, errors.New(errMsg)
	}
	return []models.History{}, nil
}

func (s *Storage) UpdateFetcher(fetcher *models.Fetcher) error {
	if s.UpdateFetcherErr {
		return errors.New(errMsg)
	}
	return nil
}

func (s *Storage) UpdateFetcherJobId(fetcherId, jobId int) error {
	if s.UpdateFetcherJobIdErr {
		return errors.New(errMsg)
	}
	return nil
}

func (s *Storage) AddHistory(history *models.History) error {
	if s.AddHistoryErr {
		return errors.New(errMsg)
	}
	return nil
}
