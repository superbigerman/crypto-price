package main

import (
	"fmt"
	"log"
	"net/http"

	"final/config"
	"final/internal/adapters/client/coindesk"
	httpadapter "final/internal/adapters/http"
	chiadapter "final/internal/adapters/http/chi"
	"final/internal/adapters/repository/postgres"
	"final/internal/usecases"
)

func main() {
	cfg := config.Load()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	repo, err := postgres.NewPriceRepositoryPostgres(connString)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	provider, err := coindesk.NewCoinDeskClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	priceService, err := usecases.NewPriceUseCase(repo, provider)
	if err != nil {
		log.Fatalf("Failed to create usecase: %v", err)
	}

	router := chiadapter.NewRouter()
	httpadapter.NewPriceHandler(priceService).RegisterRoutes(router)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
