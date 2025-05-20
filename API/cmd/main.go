package main

import (
	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	server "github.com/AmadoMuerte/BirthdayWish/API/internal/http_server"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
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
