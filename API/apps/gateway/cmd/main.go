package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/server"
	"github.com/AmadoMuerte/BirthdayWish/API/shared/pkg/logger"
)

func main() {
	runMode := flag.String("mode", "development", "Run mode: development or production")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(wd, "/../../api_local.env")
	if *runMode == "production" {
		envPath = filepath.Join(wd, "/apps/gateway/.env")
	}

	cfg, err := config.NewConfig(&envPath)
	if err != nil {
		err = fmt.Errorf("config error: %w", err)
		panic(err)
	}
	log := logger.SetupLogger(*runMode)

	authAddr := fmt.Sprintf("%s:%s", cfg.App.AuthServiceAddress, cfg.App.AuthServicePort)
	authClient, err := client.NewAuthClient(authAddr, log)
	if err != nil {
		err = fmt.Errorf("Auth client error: %s", err)
		panic(err)
	}

	wishAddr := fmt.Sprintf("%s:%s", cfg.App.WishlisterServiceAddress, cfg.App.WishlisterServicePort)
	wishClient, err := client.NewWishlisterClient(wishAddr, log)
	if err != nil {
		err = fmt.Errorf("Wishlister client error: %s", err)
		panic(err)
	}

	server := server.New(runMode, cfg, log, authClient, wishClient)
	server.Start()
}
