package boot

import (
	"os"

	"go-ocr/infrastructure/config"
	"go-ocr/infrastructure/database"
	"go-ocr/infrastructure/limiter"
	logger "go-ocr/infrastructure/log"
	"go-ocr/infrastructure/redis"
	tesseractsClient "go-ocr/infrastructure/tesseracts-client"
	"go-ocr/modules/health"
	"go-ocr/modules/ocr"
	"go-ocr/utils"

	redisThirdPartyLib "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type HandlerSetup struct {
	Limiter    *limiter.RateLimiter
	HealthHttp health.InterfaceHttp
	OcrHttp    ocr.InterfaceHttp
}

func MakeHandler() HandlerSetup {
	//initiate config
	config.Initialize()

	//initiate logger
	logger.Init(config.Conf.LogFormat, config.Conf.LogLevel)

	var err error

	//initiate a redis client
	var redisClient *redisThirdPartyLib.Client
	var redisLibInterface redis.LibInterface
	if config.Conf.Redis.EnableRedis {
		redisClient, err = redis.NewRedisClient(&config.Conf)
		if err != nil {
			log.Fatalf("failed initiate redis: %v", err)
			os.Exit(1)
		}
		//initiate a redis library interface
		redisLibInterface, err = redis.NewRedisLibInterface(redisClient)
		if err != nil {
			log.Fatalf("failed initiate redis library: %v", err)
			os.Exit(1)
		}
	}

	//setup infrastructure postgres
	var db database.HandlerDatabase
	if config.Conf.Postgres.EnablePostgres {
		db, err = database.NewDatabaseClient(&config.Conf)
		if err != nil {
			log.Fatalf("failed initiate database postgres: %v", err)
			os.Exit(1)
		}
	}

	//add limiter
	interval := utils.StringUnitToDuration(config.Conf.Interval)
	middlewareWithLimiter := limiter.NewRateLimiter(int(config.Conf.Rate), interval)

	//add tesseracts client library using gosseract
	tesseractsClientLib := tesseractsClient.NewClient()

	//health module
	var healthRepository health.RepositoryInterface
	var ocrRepository ocr.RepositoryInterface
	if config.Conf.Postgres.EnablePostgres {
		healthRepository = health.NewRepository(db.DbConn)
		ocrRepository = ocr.NewRepository(db.DbConn)
	} else {
		ocrRepository = ocr.NewInMemoryRepositoryRepositoryAdapter()
	}

	healthService := health.NewService(healthRepository, redisClient)
	healthModule := health.NewHttp(healthService)

	//ocr module
	ocrService := ocr.NewService(ocrRepository, redisLibInterface, tesseractsClientLib)
	ocrModule := ocr.NewHttp(ocrService)

	return HandlerSetup{
		Limiter:    middlewareWithLimiter,
		HealthHttp: healthModule,
		OcrHttp:    ocrModule,
	}
}
