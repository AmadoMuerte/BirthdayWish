package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/auth/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/auth/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/auth/internal/service"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/auth/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/logger"
)

func main() {
	runMode := flag.String("mode", "development", "Run mode: development or production")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(wd, "/../../api_local.env")
	if *runMode == "production" {
		envPath = filepath.Join(wd, "/apps/auth/.env")
	}

	log := logger.SetupLogger(*runMode)

	cfg, err := config.NewConfig(&envPath)
	if err != nil {
		err = fmt.Errorf("config error: %w", err)
		panic(err)
	}

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		err = fmt.Errorf("db error: %w", err)
		panic(err)
	}

	authService := service.NewAuthService(storage, log, cfg.App.SecretKey)

	server := server.New(runMode, cfg, storage, authService, log)
	server.Start()

	defer storage.Close()
}
