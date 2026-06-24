package main

import (
	"log"
	"os"
	"time"

	"final/internal/adapters/client/coindesk"
	"final/internal/adapters/repository/postgres"
	"final/internal/ports/chi"
	"final/internal/usecases"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	apiClient := coindesk.NewCoinDeskClient(
		os.Getenv("COINDESK_URL"),
		10*time.Second,
		false,
		"USD",
		os.Getenv("COINDESK_API_KEY"),
	)

	repo, err := postgres.NewPriceRepositoryPostgres(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer repo.Close()

	uc, err := usecases.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("usecase: %v", err)
	}

	chi.RunServer(uc)
}
