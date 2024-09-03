package primitive

import "time"

type Ocr struct {
	ID        int64     `gorm:"column:id"`
	ImageUrl  string    `gorm:"column:image_url"`
	Text      string    `gorm:"column:text"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

type ParameterFindOcr struct {
	Text      string
	Status    string
	PageSize  int
	Offset    int
	SortBy    string
	SortOrder string
}

type ParameterOcrHandler struct {
	Text   string
	Status string
}
