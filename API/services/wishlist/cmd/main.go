package main

import "github.com/docker/docker/api/server"

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
