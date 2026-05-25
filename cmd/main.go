package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"final/config"
	external "final/internal/adapters/coindesk"
	"final/internal/adapters/postgres"
	usecase "final/internal/usecases"
)

func main() {
	cfg := config.Load()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	// Создаём репозиторий (сам подключается к БД)
	repo, err := postgres.NewPriceRepositoryPostgres(connString)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	log.Println("✅ Connected to PostgreSQL")

	// Создаём клиент внешнего API
	apiClient, err := external.NewCoinDeskClient(cfg)
	if err != nil {
		log.Fatalf("Failed to creTE API client: %v", err)
	}

	// Создаём UseCase
	uc, err := usecase.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("Failed to create usecase: %v", err)
	}

	// HTTP handlers
	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := uc.GetPricesLast(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
