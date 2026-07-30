package app

import (
	"context"
	"log"
	"time"

	"final/internal/usecases"
	"final/pkg/config"
)

func startCron(uc usecases.PriceUseCase, cfg *config.Config) {
	go func() {
		symbols, _ := uc.GetAllSymbols(context.Background())
		if len(symbols) > 0 {
			log.Println("Cron: initial update...")
			uc.GetPricesLast(context.Background(), symbols)
		}

		interval := time.Duration(cfg.Cron.IntervalMin) * time.Minute
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Cron: updating prices...")
			symbols, _ := uc.GetAllSymbols(context.Background())
			if len(symbols) > 0 {
				uc.GetPricesLast(context.Background(), symbols)
			}
		}
	}()
}
