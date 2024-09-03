package primitive

import "time"

type OCrResponse struct {
	ID        int64     `json:"id"`
	ImageUrl  string    `json:"image_url"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HealthResponse struct {
	Db    string `json:"db"`
	Redis string `json:"redis"`
}
