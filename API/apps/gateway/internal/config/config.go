package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Mode                  string `envconfig:"RUN_MODE" default:"dev"`
		Address               string `envconfig:"APP_ADDRESS" default:"localhost"`
		Port                  string `envconfig:"APP_PORT" default:"3030"`
		SecretKey             string `envconfig:"SECRET_KEY" default:"bibibibiba"`
		AuthServicePort       string `envconfig:"AUTH_SERVICE_PORT" default:"50051"`
		WishlisterServicePort string `envconfig:"WISHLISTER_SERVICE_PORT" default:"50052"`
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
