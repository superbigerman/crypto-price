package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"final/config"
	"final/internal/adapters/client/coindesk"
	"final/internal/adapters/repository/postgres"
	portchi "final/internal/ports/chi"

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

	apiClient, err := coindesk.NewCoinDeskClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	priceUC, err := usecases.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("Failed to create usecase: %v", err)
	}

	router := portchi.NewChiRouter()

	router.Get("/prices", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := priceUC.GetPricesLast(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	router.Get("/min", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := priceUC.GetMinPrices(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	router.Get("/max", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := priceUC.GetMaxPrices(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	router.Get("/change", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		changes, err := priceUC.GetChangePercent(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(changes)
	})

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
