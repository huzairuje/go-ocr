package health

import (
	"context"
	
	"go-ocr/infrastructure/config"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CheckUpTimeDB(ctx context.Context) (err error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) CheckUpTimeDB(ctx context.Context) (err error) {
	if config.Conf.Postgres.EnablePostgres {
		db, err := r.db.WithContext(ctx).DB()
		if err != nil {
			return err
		}

		err = db.Ping()
		if err != nil {
			return err
		}
	}

	return nil
}
