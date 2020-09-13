package models

type History struct {
	FetcherId int     `json:"-"`
	Response  *string `json:"response"`
	Duration  float64 `json:"duration"`
	CreatedAt int64   `json:"created_at"`
}

func (h *History) Reset() {
	h.Response = nil
	h.Duration = 0
	h.CreatedAt = 0
}
