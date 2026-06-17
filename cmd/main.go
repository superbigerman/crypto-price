package main

import (
	"log"
	"net/http"
	"time"

	"final/internal/adapters/client/coindesk"
	"final/internal/adapters/repository/postgres"
	"final/internal/ports/chi"
	"final/internal/usecases"
)

func main() {
	connString := "postgres://macbook:postgres@localhost:5432/crypto?sslmode=disable"
	apiBaseURL := "https://min-api.cryptocompare.com"
	apiTimeout := 10 * time.Second
	apiRelaxed := true
	apiTSyms := "USD"
	apiKey := "1247506605989d7457389d3e93f5958bc2f7874b29422604dcca9cacc47e8b0d"

	repo, err := postgres.NewPriceRepositoryPostgres(connString)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	apiClient := coindesk.NewCoinDeskClient(apiBaseURL, apiTimeout, apiRelaxed, apiTSyms, apiKey)

	uc, err := usecases.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("Failed to create usecase: %v", err)
	}

	router := chi.NewChiRouter()
	handler := chi.NewPriceHandler(uc)

	router.Get("/prices", handler.GetPricesLast)
	router.Get("/min", handler.GetMinPrices)
	router.Get("/max", handler.GetMaxPrices)
	router.Get("/change", handler.GetChangePercent)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
