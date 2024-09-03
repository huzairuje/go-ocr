package health

import (
	"context"
	"errors"

	"go-ocr/infrastructure/config"
	logger "go-ocr/infrastructure/log"
	"go-ocr/modules/primitive"

	"github.com/go-redis/redis"
)

type InterfaceService interface {
	CheckUpTime(ctx context.Context) (resp primitive.HealthResponse, err error)
}

type Service struct {
	repository  RepositoryInterface
	redisClient *redis.Client
}

func NewService(repository RepositoryInterface, redisClient *redis.Client) InterfaceService {
	return &Service{
		repository:  repository,
		redisClient: redisClient,
	}
}

func (u *Service) CheckUpTime(ctx context.Context) (primitive.HealthResponse, error) {
	ctxName := "CheckUpTime"

	var postgresStatus string
	if config.Conf.Postgres.EnablePostgres {
		if u.repository == nil {
			err := errors.New("repository doesn't initiate on the boot file")
			return primitive.HealthResponse{}, err
		}

		errCheckDb := u.repository.CheckUpTimeDB(ctx)
		if errCheckDb != nil {
			logger.Error(ctx, ctxName, "got error when %s : %v", ctxName, errCheckDb)
			return primitive.HealthResponse{}, errCheckDb
		}
		postgresStatus = "healthy"
	} else {
		postgresStatus = "postgres is not enabled"
	}

	var redisStatus string
	if config.Conf.Redis.EnableRedis && u.redisClient != nil {
		errCheckRedis := u.redisClient.Ping().Err()
		if errCheckRedis != nil {
			logger.Error(ctx, ctxName, "got error when %s : %v", ctxName, errCheckRedis)
			return primitive.HealthResponse{}, errCheckRedis
		}
		redisStatus = "healthy"
	} else {
		redisStatus = "not initiated"
	}

	return primitive.HealthResponse{
		Db:    postgresStatus,
		Redis: redisStatus,
	}, nil
}
