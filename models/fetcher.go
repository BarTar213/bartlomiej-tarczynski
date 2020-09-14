package models

import (
	"errors"
	"net/url"
)

type Fetcher struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Interval int    `json:"interval"`
	JobId    int    `json:"-"`
}

func (f *Fetcher) Validate() error {
	if len(f.Url) == 0 {
		return errors.New("url can't be empty")
	}

	if f.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}

	_, err := url.ParseRequestURI(f.Url)
	if err != nil {
		return errors.New("invalid url")
	}

	return nil
}

func (f *Fetcher) Reset() {
	f.Id = 0
	f.Url = ""
	f.Interval = 0
	f.JobId = 0
}

type Id struct {
	Id int `json:"id"`
}
