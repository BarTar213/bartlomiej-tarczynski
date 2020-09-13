package models

type History struct {
	Response  string `json:"response"`
	Duration  float64 `json:"duration"`
	CreatedAt float64 `json:"created_at"`
}
