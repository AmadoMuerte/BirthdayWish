package main

import (
	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
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

	defer storage.Close()
}
