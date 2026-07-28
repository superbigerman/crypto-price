package app

import (
	"context"
	"log"
	"time"

	"final/internal/ports/chi"
	"final/internal/usecases"
	"final/pkg/config"
)

func Run(uc usecases.PriceUseCase, cfg *config.Config) {
	// Cron: обновление цен
	go func() {
		// Первый запуск сразу
		symbols, _ := uc.GetAllSymbols(context.Background())
		if len(symbols) > 0 {
			log.Println("Cron: initial update...")
			uc.GetPricesLast(context.Background(), symbols)
		}

		// Затем по расписанию
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

	// HTTP-сервер
	srv, err := chi.NewServer(uc)
	if err != nil {
		log.Fatalf("server: %v", err)
	}
	srv.Start()
}
