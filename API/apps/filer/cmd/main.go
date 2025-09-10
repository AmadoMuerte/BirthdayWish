package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/service"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/logger"
)

func main() {
	runMode := flag.String("mode", "development", "Run mode: development or production")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(wd, "/../../.env")
	if *runMode == "production" {
		envPath = filepath.Join(wd, "/apps/filer/.env")
	}

	log := logger.SetupLogger(*runMode)

	cfg, err := config.NewConfig(&envPath)
	if err != nil {
		err = fmt.Errorf("config error: %w", err)
		panic(err)
	}

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	filerService := service.NewFilerService(storage, log)

	server := server.New(runMode, cfg, storage, filerService, log)
	server.Start()
}
