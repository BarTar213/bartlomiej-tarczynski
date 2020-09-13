package storage

import "github.com/BarTar213/bartlomiej-tarczynski/models"

func (p *Postgres) GetFetchers() ([]models.Fetcher, error) {
	fetchers := make([]models.Fetcher, 0)
	err := p.db.Model(&fetchers).Select()

	return fetchers, err
}

func (p *Postgres) AddFetcher(fetcher *models.Fetcher) error {
	_, err := p.db.Model(fetcher).
		Returning("id").
		Insert()

	return err
}

func (p *Postgres) UpdateFetcher(fetcher *models.Fetcher) error {
	_, err := p.db.Model(fetcher).
		WherePK().
		Set("url=?url, interval=?interval").
		Returning("id").
		Update()

	return err
}

func (p *Postgres) DeleteFetcher(id int) error {
	_, err := p.db.Exec("DELETE FROM fetchers WHERE id=?", id)

	return err
}
