package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/logger"
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

	log := logger.SetupLogger(cfg.App.Mode)

	authAddr := fmt.Sprintf("%s:%s", cfg.App.Address, cfg.App.AuthServicePort)
	authClient, err := client.NewAuthClient(authAddr, log)
	if err != nil {
		err = fmt.Errorf("Auth client error: %s", err)
		panic(err)
	}

	server := server.New(cfg, log, authClient)
	server.Start()
}
