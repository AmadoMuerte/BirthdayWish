package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Port string `envconfig:"FILER_SERVICE_PORT" default:"50053"`
		Host string `envconfig:"FILER_SERVICE_HOST" default:"0.0.0.0"`
	}
	DB struct {
		Host       string        `envconfig:"FILER_DB_HOST" default:"localhost"`
		Name       string        `envconfig:"FILER_DB_NAME" default:"filer_service"`
		User       string        `envconfig:"FILER_DB_USER" default:"minioadmin"`
		Pass       string        `envconfig:"FILER_DB_PASS" default:"minioadmin"`
		Port       string        `envconfig:"FILER_DB_PORT" default:"9000"`
		BucketName string        `envconfig:"FILER_DB_BUCKET_NAME" default:"imgs"`
		UseSSL     bool          `envconfig:"FILER_DB_USE_SSL" default:"false"`
		Timeout    time.Duration `envconfig:"FILER_DB_TIMEOUT" default:"5m"`
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
