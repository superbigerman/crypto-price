package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"final/internal/adapters/client/coindesk"
	"final/internal/adapters/repository/postgres"
	"final/internal/usecases"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// ========== Шаг 1: Создаём настоящие зависимости ==========
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

	// ========== Шаг 2: Трейсинг — ручные данные ==========
	ctx := context.Background()
	testSymbols := []string{"BTC", "ETH"}

	fmt.Println("=== ТРЕЙСИНГ: GetPricesLast ===")
	fmt.Printf("Входные данные: symbols=%v\n", testSymbols)

	prices, err := uc.GetPricesLast(ctx, testSymbols)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Результат: %+v\n", prices)
	}

	fmt.Println("\n=== ТРЕЙСИНГ: GetMinPrices ===")
	fmt.Printf("Входные данные: symbols=%v\n", testSymbols)

	prices, err = uc.GetMinPrices(ctx, testSymbols)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Результат: %+v\n", prices)
	}

	fmt.Println("\n=== ТРЕЙСИНГ: GetMaxPrices ===")
	fmt.Printf("Входные данные: symbols=%v\n", testSymbols)

	prices, err = uc.GetMaxPrices(ctx, testSymbols)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Результат: %+v\n", prices)
	}

	fmt.Println("\n=== ТРЕЙСИНГ: GetChangePercent ===")
	fmt.Printf("Входные данные: symbols=%v\n", testSymbols)

	changes, err := uc.GetChangePercent(ctx, testSymbols)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
	} else {
		fmt.Printf("✅ Результат: %+v\n", changes)
	}

	fmt.Println("\n=== ТРЕЙСИНГ ЗАВЕРШЁН ===")
}
