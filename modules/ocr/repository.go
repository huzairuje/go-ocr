package ocr

import (
	"context"
	"strings"

	"go-ocr/modules/primitive"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateOcr(ctx context.Context, request primitive.Ocr) (result primitive.Ocr, err error)
	FindOcrByID(ctx context.Context, id int64) (result primitive.Ocr, err error)
	FindOcrByText(ctx context.Context, text string) (result primitive.Ocr, err error)
	FindAllListOcrPagination(ctx context.Context, param primitive.ParameterFindOcr) (result []primitive.Ocr, err error)
	CountAllListOcr(ctx context.Context, param primitive.ParameterFindOcr) (count int64, err error)
	FindAllListOcrNonPagination(ctx context.Context, param primitive.ParameterFindOcr) (result []primitive.Ocr, err error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CreateOcr(ctx context.Context, request primitive.Ocr) (result primitive.Ocr, err error) {
	err = repo.db.WithContext(ctx).Table("ocr").Create(&request).Scan(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repo *Repository) FindOcrByID(ctx context.Context, id int64) (result primitive.Ocr, err error) {
	err = repo.db.WithContext(ctx).Table("ocr").
		Where("id = ?", id).
		Where("deleted_at is null").
		First(&result).
		Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repo *Repository) FindOcrByText(ctx context.Context, text string) (result primitive.Ocr, err error) {
	err = repo.db.WithContext(ctx).Table("ocr").
		Where("text ilike '%?%'", text).
		Where("deleted_at is null").
		First(&result).
		Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repo *Repository) FindAllListOcrPagination(ctx context.Context, param primitive.ParameterFindOcr) (result []primitive.Ocr, err error) {
	query := repo.db.WithContext(ctx).Table("ocr")

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if param.Text != "" {
		query = query.Where("text ilike '%?%'", param.Text)
	}

	err = query.Offset(param.Offset).
		Limit(param.PageSize).
		Order(strings.Join([]string{param.SortBy, param.SortOrder}, " ")).
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *Repository) CountAllListOcr(ctx context.Context, param primitive.ParameterFindOcr) (count int64, err error) {
	query := repo.db.WithContext(ctx).Table("ocr")

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if param.Text != "" {
		query = query.Where("text ilike '%?%'", param.Text)
	}

	err = query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *Repository) FindAllListOcrNonPagination(ctx context.Context, param primitive.ParameterFindOcr) (result []primitive.Ocr, err error) {
	query := repo.db.WithContext(ctx).Table("ocr")

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if param.Text != "" {
		query = query.Where("text ilike '%?%'", param.Text)
	}

	err = query.Order("id desc").
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
