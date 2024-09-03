package ocr

import "gorm.io/gorm"

type RepositoryInterface interface {
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
