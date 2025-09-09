package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Port string `envconfig:"WISHLISTER_SERVICE_PORT" default:"50052"`
		Host string `envconfig:"APP_ADDRESS" default:"0.0.0.0"`
		Mode string `envconfig:"RUN_MODE" default:"dev"`
	}
	DB struct {
		Host string `envconfig:"DB_HOST" default:"localhost"`
		Name string `envconfig:"WISH_DB_NAME" default:"wish_service"`
		User string `envconfig:"WISH_DB_USER" default:"postgres"`
		Pass string `envconfig:"WISH_DB_PASS" default:"postgres"`
		Port string `envconfig:"WISH_DB_PORT" default:"5434"`
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
