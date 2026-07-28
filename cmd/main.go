package main

import (
	"log"
	"time"

	"final/internal/adapters/client/coindesk"
	"final/internal/adapters/repository/postgres"
	"final/internal/app"
	"final/internal/usecases"
	"final/pkg/config"
)

// @title Crypto Price API
// @version 1.0
// @description Сервис для получения курсов криптовалют
// @host localhost:8080
// @BasePath /

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	apiClient := coindesk.NewCoinDeskClient(
		cfg.CoinDesk.URL,
		time.Duration(cfg.CoinDesk.TimeoutSec)*time.Second,
		false,
		"USD",
		cfg.CoinDesk.APIKey,
	)

	repo, err := postgres.NewPriceRepositoryPostgres(cfg.Database.URL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer repo.Close()

	uc, err := usecases.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("usecase: %v", err)
	}

	app.Run(uc, cfg)
}
