package main

import (
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
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	runMode := os.Args[1]
	envPath := filepath.Join(wd, "/../../.env")
	if runMode == "production" {
		envPath = filepath.Join(wd, "/apps/auth/.env")
	}

	cfg, err := config.NewConfig(&envPath)
	if err != nil {
		err = fmt.Errorf("config error: %w", err)
		panic(err)
	}

	log := logger.SetupLogger(runMode)

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		err = fmt.Errorf("db error: %w", err)
		panic(err)
	}

	authService := service.NewAuthService(storage, log, cfg.App.SecretKey)

	server := server.New(cfg, storage, authService, log)
	server.Start()

	defer storage.Close()
}
