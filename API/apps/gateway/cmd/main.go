package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/config"
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
		err = fmt.Errorf("Config error: %s", err)
		panic(err)
	}
	storage, err := storage.NewStorage(cfg)
	if err != nil {
		err = fmt.Errorf("DB error: %s", err)
		panic(err)
	}

	rdb, err := redis.GetInstance(cfg)
	if err != nil {
		panic(err)
	}
	defer rdb.Client.Close()

	server := server.New(cfg, storage, rdb)
	server.Start()

	defer storage.Close()
}
