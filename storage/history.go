package storage

import "github.com/BarTar213/bartlomiej-tarczynski/models"

func (p *Postgres) GetHistory(id int) ([]models.History, error) {
	history := make([]models.History, 0)
	err := p.db.Model(&history).
		Where("fetcher_id=?", id).
		Select()

	return history, err
}

func (p *Postgres) AddHistory(history *models.History) error {
	_, err := p.db.Model(&history).Insert()

	return err
}
