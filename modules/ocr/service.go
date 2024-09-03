package ocr

import (
	"github.com/go-redis/redis"
)

type ServiceInterface interface {
}

type Service struct {
	repository  RepositoryInterface
	redisClient *redis.Client
}

func NewService(repository RepositoryInterface, redisClient *redis.Client) ServiceInterface {
	return &Service{
		repository:  repository,
		redisClient: redisClient,
	}
}
