package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/worker"
	mainconfig "github.com/AmadoMuerte/BirthdayWish/API/pkg/config"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/redis"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(wd, "/../../.env")
	cfg, err := config.NewConfig(&envPath)
	if err != nil {
		panic(err)
	}

	mainCfg, err := mainconfig.NewConfig(&envPath)
	if err != nil {
		panic(err)
	}

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	rdb, err := redis.GetInstance(mainCfg)
	if err != nil {
		panic(err)
	}
	defer rdb.Client.Close()

	worker := worker.New(rdb, slog.Default(), storage)
	go worker.Start(context.Background())

	server := server.New(cfg, storage)
	server.Start()

	defer storage.Close()
}
