package ocr

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"go-ocr/modules/primitive"
)

// InMemoryRepository stores OCR data in memory.
type InMemoryRepository struct {
	ocrs       []primitive.Ocr
	idSequence int64
	mu         sync.RWMutex
}

// CreateOcr adds a new OCR entry to the in-memory repository.
func (i *InMemoryRepository) CreateOcr(ctx context.Context, request primitive.Ocr) (result primitive.Ocr, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Assign a new ID from the sequence and increment it.
	request.ID = i.idSequence
	i.idSequence++

	// Add the OCR entry to the repository.
	i.ocrs = append(i.ocrs, request)

	// Return the newly created OCR entry.
	return request, nil
}

// FindOcrByID retrieves an OCR entry by its ID.
func (i *InMemoryRepository) FindOcrByID(ctx context.Context, id int64) (result primitive.Ocr, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Search for the OCR entry with the matching ID.
	for _, ocr := range i.ocrs {
		if ocr.ID == id {
			return ocr, nil
		}
	}

	// Return an error if not found.
	return primitive.Ocr{}, errors.New("OCR entry not found")
}

// FindOcrByText retrieves an OCR entry by matching its text.
func (i *InMemoryRepository) FindOcrByText(ctx context.Context, text string) (result primitive.Ocr, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Search for the OCR entry with the matching text.
	for _, ocr := range i.ocrs {
		if ocr.Text == text {
			return ocr, nil
		}
	}

	// Return an error if not found.
	return primitive.Ocr{}, errors.New("OCR entry not found")
}

// FindAllListOcrPagination returns a paginated and filtered list of OCR entries.
func (i *InMemoryRepository) FindAllListOcrPagination(ctx context.Context, param primitive.ParameterFindOcr) (result []primitive.Ocr, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Filter based on Text and Status
	filtered := make([]primitive.Ocr, 0)
	for _, ocr := range i.ocrs {
		if (param.Text == "" || strings.Contains(ocr.Text, param.Text)) &&
			(param.Status == "" || ocr.Status == param.Status) && ocr.DeletedAt.IsZero() {
			filtered = append(filtered, ocr)
		}
	}

	// Sort by the requested field and order
	if param.SortBy != "" {
		switch param.SortBy {
		case "Text":
			if param.SortOrder == "asc" {
				sort.Slice(filtered, func(i, j int) bool {
					return filtered[i].Text < filtered[j].Text
				})
			} else {
				sort.Slice(filtered, func(i, j int) bool {
					return filtered[i].Text > filtered[j].Text
				})
			}
		case "ID":
			if param.SortOrder == "asc" {
				sort.Slice(filtered, func(i, j int) bool {
					return filtered[i].ID < filtered[j].ID
				})
			} else {
				sort.Slice(filtered, func(i, j int) bool {
					return filtered[i].ID > filtered[j].ID
				})
			}
		}
	}

	// Apply pagination
	start := param.Offset
	end := start + param.PageSize

	if start > len(filtered) {
		return []primitive.Ocr{}, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	// Return the paginated result
	return filtered[start:end], nil
}

// CountAllListOcr counts the total number of OCR entries in the repository.
func (i *InMemoryRepository) CountAllListOcr(ctx context.Context, param primitive.ParameterFindOcr) (count int64, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Filter based on Text and Status
	filtered := make([]primitive.Ocr, 0)
	for _, ocr := range i.ocrs {
		if (param.Text == "" || strings.Contains(ocr.Text, param.Text)) && (param.Status == "" || ocr.Status == param.Status) && ocr.DeletedAt.IsZero() {
			filtered = append(filtered, ocr)
		}
	}

	// Return the count of filtered results.
	return int64(len(filtered)), nil
}

// FindAllListOcrNonPagination returns all OCR entries without pagination.
func (i *InMemoryRepository) FindAllListOcrNonPagination(ctx context.Context, param primitive.ParameterFindOcr) (result []primitive.Ocr, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Filter based on Text and Status
	filtered := make([]primitive.Ocr, 0)
	for _, ocr := range i.ocrs {
		if (param.Text == "" || strings.Contains(ocr.Text, param.Text)) &&
			(param.Status == "" || ocr.Status == param.Status) && ocr.DeletedAt.IsZero() {
			filtered = append(filtered, ocr)
		}
	}

	// Return the filtered list without pagination.
	return filtered, nil
}

// NewInMemoryRepository creates a new instance of InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		ocrs:       make([]primitive.Ocr, 0),
		idSequence: 1,
	}
}

// NewInMemoryRepositoryRepositoryAdapter creates a new instance of RepositoryInterface using InMemoryRepository.
func NewInMemoryRepositoryRepositoryAdapter() RepositoryInterface {
	return NewInMemoryRepository()
}
