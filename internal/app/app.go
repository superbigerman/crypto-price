package app

import (
	"fmt"
	"log"
	"time"

	"final/deploy/config"
	"final/internal/adapters/client/coindesk"
	"final/internal/adapters/repository/postgres"
	"final/internal/ports/chi"
	"final/internal/usecases"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	repo, err := postgres.NewPriceRepositoryPostgres(cfg.Database.URL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer repo.Close()
	apiClient := coindesk.NewCoinDeskClient(
		cfg.CoinDesk.URL,
		time.Duration(cfg.CoinDesk.TimeoutSec)*time.Second,
		false,
		"USD",
		cfg.CoinDesk.APIKey,
	)
	uc, err := usecases.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("usecase: %v", err)
	}

	startCron(uc, cfg)

	srv, err := chi.NewServer(uc)
	if err != nil {
		log.Fatalf("server: %v", err)
	}
	srv.Start()
	return nil
}
