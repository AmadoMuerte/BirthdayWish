package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type MinIOConfig struct {
	Endpoint    string        `envconfig:"MINIO_ENDPOINT" default:"localhost:9000"`
	AccessKey   string        `envconfig:"MINIO_USER" default:"admin"`
	SecretKey   string        `envconfig:"MINIO_PASS" default:"StrongPass123"`
	BucketName  string        `envconfig:"MINIO_BUCKET_NAME" default:"imgs"`
	UseSSL      bool          `envconfig:"MINIO_USE_SSL" default:"false"`
	APIPort     string        `envconfig:"MINIO_API_PORT" default:"9000"`
	ConsolePort string        `envconfig:"MINIO_CONSOLE_PORT" default:"9001"`
	Timeout     time.Duration `envconfig:"MINIO_TIMEOUT" default:"5m"`
}

type AppConfig struct {
	Mode    string `envconfig:"RUN_MODE" default:"dev"`
	Address string `envconfig:"MINLO_ADDRESS" default:"localhost:5058"`
}

type Config struct {
	MinIO MinIOConfig
	App   AppConfig
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
