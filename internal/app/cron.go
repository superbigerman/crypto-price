package app

import (
	"context"
	"fmt"
	"log"

	"final/deploy/config"
	"final/internal/usecases"

	"github.com/robfig/cron/v3"
)

func startCron(uc usecases.PriceUseCase, cfg *config.Config) {
	c := cron.New()

	expression := fmt.Sprintf("*/%d * * * *", cfg.Cron.IntervalMin)
	c.AddFunc(expression, func() {
		log.Println("Cron: updating prices...")
		if err := uc.UpdateAllPrices(context.Background()); err != nil {
			log.Printf("Cron error: %v", err)
		}
	})

	// Первый запуск сразу
	go func() {
		log.Println("Cron: initial update...")
		if err := uc.UpdateAllPrices(context.Background()); err != nil {
			log.Printf("Cron error: %v", err)
		}
	}()

	c.Start()
}
