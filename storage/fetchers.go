package storage

import (
	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/go-pg/pg/v10"
)

func (p *Postgres) GetFetchers() ([]models.Fetcher, error) {
	fetchers := make([]models.Fetcher, 0)
	err := p.db.Model(&fetchers).Select()

	return fetchers, err
}

func (p *Postgres) GetFetchersForSync() ([]models.Fetcher, error) {
	fetchers := make([]models.Fetcher, 0)
	err := p.db.Model(&fetchers).Where("job_id=0").Select()

	return fetchers, err
}

func (p *Postgres) GetFetcherJob(id int) (int, error){
	var jobId int
	_, err := p.db.QueryOne(pg.Scan(&jobId), "SELECT job_id FROM fetchers WHERE id=?", id)

	return jobId, err
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
		Returning("id, job_id").
		Update()

	return err
}

func (p *Postgres) UpdateFetcherJobId(fetcherId, jobId int) error {
	_, err := p.db.Exec("UPDATE fetchers SET job_id=? WHERE id=?", fetcherId, jobId)

	return err
}

func (p *Postgres) DeleteFetcher(id int) (int, error) {
	var jobId int
	_, err := p.db.QueryOne(pg.Scan(&jobId), "DELETE FROM fetchers WHERE id=? RETURNING job_id", id)

	return jobId, err
}
