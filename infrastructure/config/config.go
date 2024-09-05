package config

import (
	"errors"
	"os"
	"strings"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func initialiseRemote(v *viper.Viper) error {
	consulUrl := os.Getenv("CONSUL_URL")
	_ = v.AddRemoteProvider("consul", consulUrl, "GO_OCR")
	v.SetConfigType("yaml")
	return v.ReadRemoteConfig()
}

func initialiseFileAndEnv(v *viper.Viper, env string) error {
	v.SetConfigName(configName[env])
	for _, path := range searchPath {
		v.AddConfigPath(path)
	}
	v.SetEnvPrefix("GO_OCR")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return v.ReadInConfig()
}

func initialiseDefaults(v *viper.Viper) {
	for key, value := range configDefaults {
		v.SetDefault(key, value)
	}
}

func Initialize() {
	v := viper.New()
	initialiseDefaults(v)
	if err := initialiseRemote(v); err != nil {
		log.Warningf("No remote server configured will load configuration from file and environment variables: %+v", err)
		if err := initialiseFileAndEnv(v, Env); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				log.Warning("No 'config.yaml' file found on search paths. Will either use environment variables or defaults")
			}
		}
	}

	err := v.Unmarshal(&Conf)
	if err != nil {
		log.Printf("Error un-marshalling configuration: %s", err.Error())
	}
}

var (
	Conf        Config
	Env         string
	RedisClient *redis.Client

	EnvironmentLocal = "LOCAL"
	EnvironmentDev   = "DEV"
	EnvironmentUAT   = "UAT"
	EnvironmentProd  = "PROD"
	ListOfIsland     map[uint64]string

	searchPath = []string{
		"/etc/go-ocr",
		"$HOME/.go-ocr",
		".",
	}
	configDefaults = map[string]interface{}{
		"port":       1234,
		"logLevel":   "DEBUG",
		"logFormat":  "text",
		"signString": "supersecret",
	}
	configName = map[string]string{
		"local": "config.local",
		"dev":   "config.dev",
		"uat":   "config.uat",
		"prod":  "config.prod",
		"test":  "config.test",
	}
)

type Config struct {
	Env              string           `mapstructure:"env"`
	Port             int              `mapstructure:"port"`
	LogLevel         string           `mapstructure:"logLevel"`
	LogMode          bool             `mapstructure:"logMode"`
	LogFormat        string           `mapstructure:"logFormat"`
	Postgres         PostgresConfig   `mapstructure:"postgres"`
	Redis            RedisConfig      `mapstructure:"redis"`
	Rate             int64            `mapstructure:"rate"`
	Interval         string           `mapstructure:"interval"`
	TesseractsConfig TesseractsConfig `mapstructure:"tesseracts"`
}

// PostgresConfig ...
type PostgresConfig struct {
	ConnMaxLifetime    int    `mapstructure:"connectTimeout"`
	MaxOpenConnections int    `mapstructure:"maxOpenConnections"`
	MaxIdleConnections int    `mapstructure:"maxIdleConnections"`
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	Schema             string `mapstructure:"schema"`
	DBName             string `mapstructure:"dbName"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	EnablePostgres     bool   `mapstructure:"enablePostgres"`
}

type RedisConfig struct {
	Host        string `mapstructure:"host"`
	Password    string `mapstructure:"password"`
	DB          int    `mapstructure:"db"`
	Port        int    `mapstructure:"port"`
	EnableRedis bool   `mapstructure:"enableRedis"`
}

type TesseractsConfig struct {
	Languages []string `mapstructure:"languages"`
}
