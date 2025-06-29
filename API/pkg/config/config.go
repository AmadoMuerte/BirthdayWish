package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type db struct {
	Host string `required:"true" envconfig:"DB_HOST"`
	Name string `required:"true" envconfig:"DB_NAME"`
	User string `required:"true" envconfig:"DB_USER"`
	Pass string `required:"true" envconfig:"DB_PASS"`
	Port string `required:"true" envconfig:"DB_PORT"`
}

type redis struct {
	Path string `required:"true" envconfig:"REDIS_PATH"`
	Port string `required:"true" envconfig:"REDIS_PORT"`
	Pass string `required:"true" envconfig:"REDIS_PASS"`
	User string `required:"true" envconfig:"REDIS_USER"`
}

type appConfig struct {
	Mode      string `envconfig:"RUN_MODE" default:"dev"`
	Address   string `envconfig:"APP_ADDRESS" default:"localhost"`
	Port      string `envconfig:"APP_PORT" default:"8080"`
	SecretKey string `envconfig:"SECRET_KEY"`
}

type services struct {
	WishListAddr string `envconfig:"WISHLIST_ADDRESS" default:"localhost:5057"`
	Minio        string `envconfig:"MINLO_ADDRESS" default:"localhost:5058"`
}

type Config struct {
	DB       db
	Redis    redis
	App      appConfig
	Services services
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
