package ocr

import (
	"go-ocr/modules/primitive"
	"sync"
)

type InMemoryRepository struct {
	articles   []primitive.Ocr
	idSequence int64
	mu         sync.RWMutex
}

// NewInMemoryRepository creates a new instance of InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		articles:   make([]primitive.Ocr, 0),
		idSequence: 1,
	}
}

// NewInMemoryRepositoryRepositoryAdapter creates a new instance of RepositoryInterface using InMemoryRepository.
func NewInMemoryRepositoryRepositoryAdapter() RepositoryInterface {
	return NewInMemoryRepository()
}
