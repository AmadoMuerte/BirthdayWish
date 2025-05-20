package main

import (
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/config"
	server "github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/http_server"
	"github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/storage"
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

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	server := server.New(cfg, storage)
	server.Start()

	defer storage.Close()
}
