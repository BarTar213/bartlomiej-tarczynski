package models

type Fetcher struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Interval int    `json:"interval"`
}
