package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go-ocr/infrastructure/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	defaultConnMaxLifeTime int = 5  // default max 5 minutes lifetime
	defaultMaxOpenConns    int = 10 // default max 10 open connections
	defaultMaxIdleConns    int = 10 // default max 10 idle connections
)

type HandlerDatabase struct {
	DbConn *gorm.DB
}

func NewDatabaseClient(conf *config.Config) (HandlerDatabase, error) {
	db := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s fallback_application_name=go-soskomlap-api-master TimeZone=Asia/Jakarta",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.DBName,
		"disable")
	dbConn, err := loadPsqlDb(conf.LogMode, db, conf.Postgres.MaxOpenConnections, conf.Postgres.MaxIdleConnections, conf.Postgres.ConnMaxLifetime)
	if err != nil {
		log.Printf("failed to connect database instance: %v", err)
		return HandlerDatabase{}, err
	}

	return HandlerDatabase{
		DbConn: dbConn,
	}, nil
}

func loadPsqlDb(logMode bool, psqlInfo string, maxOpenConn, maxIdleConn, maxLifetime int) (*gorm.DB, error) {
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// checking if connection to db has been established
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	if maxLifetime == 0 {
		maxLifetime = defaultConnMaxLifeTime
	}

	if maxIdleConn == 0 {
		maxIdleConn = defaultMaxIdleConns
	}

	if maxOpenConn == 0 {
		maxOpenConn = defaultMaxOpenConns
	}

	conn.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)
	conn.SetMaxOpenConns(maxOpenConn)
	conn.SetMaxIdleConns(maxIdleConn)

	var gormConfig *gorm.Config
	if logMode {
		gormConfig = &gorm.Config{
			PrepareStmt: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		gormConfig = &gorm.Config{
			PrepareStmt: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: nil,
		}
	}

	return gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), gormConfig)
}
