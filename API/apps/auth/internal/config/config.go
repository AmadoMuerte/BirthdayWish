package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Port      string `envconfig:"AUTH_SERVICE_PORT" default:"50051"`
		Host      string `envconfig:"AUTH_SERVICE_HOST" default:"0.0.0.0"`
		SecretKey string `envconfig:"SECRET_KEY" default:"bibibibiba"`
	}
	DB struct {
		Host string `envconfig:"AUTH_DB_HOST" default:"localhost"`
		Name string `envconfig:"AUTH_DB_NAME" default:"users_service"`
		User string `envconfig:"AUTH_DB_USER" default:"postgres"`
		Pass string `envconfig:"AUTH_DB_PASS" default:"postgres"`
		Port string `envconfig:"AUTH_DB_PORT" default:"5433"`
	}
}

var getWd = os.Getwd
var processEnv = envconfig.Process

func NewConfig(customPath *string) (*Config, error) {
	var newCfg Config

	wd, err := getWd()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(wd, ".env")

	if customPath != nil {
		envPath = *customPath
	}

	_ = godotenv.Overload(envPath)
	if err = processEnv("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
}
