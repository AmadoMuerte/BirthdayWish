package main

import (
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/services/minio/internal/config"
	server "github.com/AmadoMuerte/BirthdayWish/API/services/minio/internal/http_server"
	"github.com/AmadoMuerte/BirthdayWish/API/services/minio/internal/storage"
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

}
