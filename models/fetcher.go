package models

import "errors"

type Fetcher struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Interval int    `json:"interval"`
}

func (f *Fetcher) Validate() error {
	if len(f.Url) == 0 {
		return errors.New("url can't be empty")
	}

	if f.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}

	return nil
}
