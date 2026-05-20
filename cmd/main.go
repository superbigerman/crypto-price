package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"final/internal/config"
	"final/internal/external"
	"final/internal/repository/postgres"
	"final/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// Строка подключения к PostgreSQL
	connString := postgres.BuildConnString(
		cfg.DBHost, cfg.DBPort, cfg.DBUser,
		cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	// Подключаемся к базе данных
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to PostgreSQL")

	// Создаём репозиторий
	repo := postgres.NewPriceRepositoryPostgres(pool)

	// Создаём клиент внешнего API
	apiClient := external.NewCoinDeskClient(cfg)

	// Создаём UseCase (с проверкой ошибки)
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

	http.HandleFunc("/min", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := uc.GetMinPrices(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	http.HandleFunc("/max", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := uc.GetMaxPrices(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		changes, err := uc.GetChangePercent(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(changes)
	})

	// Запускаем сервер
	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
