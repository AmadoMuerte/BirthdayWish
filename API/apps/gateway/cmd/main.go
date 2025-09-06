package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/logger"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	envPath := filepath.Join(wd, "/../../.env")

	cfg, err := config.NewConfig(&envPath)
	if err != nil {
		err = fmt.Errorf("Config error: %s", err)
		panic(err)
	}
	storage, err := storage.NewStorage(cfg)
	if err != nil {
		err = fmt.Errorf("DB error: %s", err)
		panic(err)
	}

	log := logger.SetupLogger(cfg.App.Mode)

	server := server.New(cfg, storage, log)
	server.Start()

	defer storage.Close()
}
