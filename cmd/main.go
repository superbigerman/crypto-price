package main

import (
	"encoding/json"
	"final/config"
	"final/internal/adapters/chi"
	"final/internal/adapters/coindesk"
	"final/internal/adapters/postgres"
	"final/internal/ports"
	"final/internal/usecases"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	cfg := config.Load()

	// PostgreSQL
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	repo, err := postgres.NewPriceRepositoryPostgres(connString)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// CoinDesk клиент
	apiClient, err := coindesk.NewCoinDeskClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	// UseCase
	uc, err := usecases.NewPriceUseCase(repo, apiClient)
	if err != nil {
		log.Fatalf("Failed to create usecase: %v", err)
	}

	// Chi роутер
	router := chi.NewChiRouter()

	router.Group(func(r ports.Router) {
		r.Get("/prices", func(w http.ResponseWriter, req *http.Request) {
			symbolsParam := req.URL.Query().Get("symbols")
			if symbolsParam == "" {
				http.Error(w, "missing symbols param", http.StatusBadRequest)
				return
			}
			symbols := strings.Split(symbolsParam, ",")

			prices, err := uc.GetPricesLast(req.Context(), symbols)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(prices)
		})

		r.Get("/min", func(w http.ResponseWriter, req *http.Request) {
			symbolsParam := req.URL.Query().Get("symbols")
			if symbolsParam == "" {
				http.Error(w, "missing symbols param", http.StatusBadRequest)
				return
			}
			symbols := strings.Split(symbolsParam, ",")

			prices, err := uc.GetMinPrices(req.Context(), symbols)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(prices)
		})

		r.Get("/max", func(w http.ResponseWriter, req *http.Request) {
			symbolsParam := req.URL.Query().Get("symbols")
			if symbolsParam == "" {
				http.Error(w, "missing symbols param", http.StatusBadRequest)
				return
			}
			symbols := strings.Split(symbolsParam, ",")

			prices, err := uc.GetMaxPrices(req.Context(), symbols)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(prices)
		})

		r.Get("/change", func(w http.ResponseWriter, req *http.Request) {
			symbolsParam := req.URL.Query().Get("symbols")
			if symbolsParam == "" {
				http.Error(w, "missing symbols param", http.StatusBadRequest)
				return
			}
			symbols := strings.Split(symbolsParam, ",")

			changes, err := uc.GetChangePercent(req.Context(), symbols)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(changes)
		})
	})

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
