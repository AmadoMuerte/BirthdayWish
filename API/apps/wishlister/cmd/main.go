package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/service"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
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
		envPath = filepath.Join(wd, "/apps/wishlister/.env")
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

	wishService := service.NewWishService(storage, log)

	server := server.New(cfg, storage, wishService, log)
	server.Start()

	defer storage.Close()
}
