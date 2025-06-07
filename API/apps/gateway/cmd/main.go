package main

import (
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
)

func main() {
	cfg, err := config.NewConfig(nil)
	if err != nil {
		panic(err)
	}
	storage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	server := server.New(cfg, storage)
	server.Start()

	defer storage.Close()
}
