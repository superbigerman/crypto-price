package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"final/internal/entity"
	"final/internal/usecase"
)

// ========== МОК-РЕПОЗИТОРИЙ ==========

type mockRepo struct {
	prices map[string][]entity.Price
}

func (m *mockRepo) GetPrice(symbol string) (entity.Price, error) {
	if m.prices == nil {
		return entity.Price{}, fmt.Errorf("repository not initialized")
	}
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return entity.Price{}, fmt.Errorf("price for %s not found", symbol)
	}
	return prices[len(prices)-1], nil
}

func (m *mockRepo) SavePrice(price entity.Price) error {
	if m.prices == nil {
		m.prices = make(map[string][]entity.Price)
	}
	m.prices[price.Symbol] = append(m.prices[price.Symbol], price)
	return nil
}

func (m *mockRepo) GetMinPrice(symbol string) (entity.MinPriceResponse, error) {
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return entity.MinPriceResponse{}, fmt.Errorf("no prices for %s", symbol)
	}

	minPrice := prices[0].Price
	minTime := prices[0].CreatedAt
	for _, p := range prices {
		if p.Price < minPrice {
			minPrice = p.Price
			minTime = p.CreatedAt
		}
	}

	return *entity.NewMinPriceResponse(symbol, minPrice, minTime), nil
}

func (m *mockRepo) GetMaxPrice(symbol string) (entity.MaxPriceResponse, error) {
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return entity.MaxPriceResponse{}, fmt.Errorf("no prices for %s", symbol)
	}

	maxPrice := prices[0].Price
	maxTime := prices[0].CreatedAt
	for _, p := range prices {
		if p.Price > maxPrice {
			maxPrice = p.Price
			maxTime = p.CreatedAt
		}
	}

	return *entity.NewMaxPriceResponse(symbol, maxPrice, maxTime), nil
}

func (m *mockRepo) GetChangePercent(symbol string) (entity.ChangePercentResponse, error) {
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return entity.ChangePercentResponse{}, fmt.Errorf("no prices for %s", symbol)
	}

	currentPrice := prices[len(prices)-1].Price
	hourAgo := time.Now().Add(-1 * time.Hour)
	var hourAgoPrice float64
	found := false

	for i := len(prices) - 1; i >= 0; i-- {
		if prices[i].CreatedAt.Before(hourAgo) || prices[i].CreatedAt.Equal(hourAgo) {
			hourAgoPrice = prices[i].Price
			found = true
			break
		}
	}

	if !found {
		return entity.ChangePercentResponse{}, fmt.Errorf("no price data for 1 hour ago")
	}

	if hourAgoPrice == 0 {
		return entity.ChangePercentResponse{}, fmt.Errorf("cannot calculate change: price hour ago was zero")
	}

	changePercent := ((currentPrice - hourAgoPrice) / hourAgoPrice) * 100

	direction := "stable"
	if changePercent > 0 {
		direction = "up"
	} else if changePercent < 0 {
		direction = "down"
	}

	return *entity.NewChangePercentResponse(symbol, changePercent, direction), nil
}

func (m *mockRepo) GetPrices(symbols []string) (map[string]entity.Price, error) {
	if m.prices == nil {
		return nil, fmt.Errorf("repository not initialized")
	}

	result := make(map[string]entity.Price)
	for _, symbol := range symbols {
		prices, ok := m.prices[symbol]
		if !ok || len(prices) == 0 {
			continue
		}
		result[symbol] = prices[len(prices)-1]
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no prices found for requested symbols")
	}

	return result, nil
}

// ========== МОК-ВНЕШНИЙ API ==========

type mockExternalAPI struct{}

func (m *mockExternalAPI) GetRealTimePrice(symbol string) (entity.Price, error) {
	mockPrices := map[string]float64{
		"BTC":  45000,
		"ETH":  3200,
		"XRP":  0.5,
		"DOGE": 0.08,
		"SOL":  100,
	}
	price, ok := mockPrices[symbol]
	if !ok {
		return entity.Price{}, fmt.Errorf("price for %s not available", symbol)
	}

	newPrice, err := entity.NewPrice(symbol, price) // ✅ получаем оба значения
	if err != nil {
		return entity.Price{}, err
	}
	return *newPrice, nil
}

// ========== MAIN ==========

func main() {
	repo := &mockRepo{}
	externalAPI := &mockExternalAPI{}
	uc := usecase.NewPriceUseCase(repo, externalAPI)

	// Тестовые данные для BTC
	btcPrice1, _ := entity.NewPrice("BTC", 44000)
	btcPrice2, _ := entity.NewPrice("BTC", 44500)
	btcPrice3, _ := entity.NewPrice("BTC", 45000)
	btcPrice4, _ := entity.NewPrice("BTC", 44800)
	repo.SavePrice(*btcPrice1)
	time.Sleep(2 * time.Second)
	repo.SavePrice(*btcPrice2)
	time.Sleep(2 * time.Second)
	repo.SavePrice(*btcPrice3)
	time.Sleep(2 * time.Second)
	repo.SavePrice(*btcPrice4)

	// Тестовые данные для ETH
	ethPrice1, _ := entity.NewPrice("ETH", 3100)
	ethPrice2, _ := entity.NewPrice("ETH", 3200)
	ethPrice3, _ := entity.NewPrice("ETH", 3250)
	repo.SavePrice(*ethPrice1)
	time.Sleep(2 * time.Second)
	repo.SavePrice(*ethPrice2)
	time.Sleep(2 * time.Second)
	repo.SavePrice(*ethPrice3)

	// Тестовые данные для XRP
	xrpPrice, _ := entity.NewPrice("XRP", 0.5)
	repo.SavePrice(*xrpPrice)

	fmt.Println("✅ Тестовые данные загружены")

	// ========== HTTP HANDLERS ==========

	// GET /price/{symbol}
	http.HandleFunc("/price/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/price/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		price, err := uc.GetPrice(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol": price.Symbol,
			"price":  price.Price,
		})
	})

	// GET /min/{symbol}
	http.HandleFunc("/min/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/min/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		response, err := uc.GetMinPrice(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// GET /max/{symbol}
	http.HandleFunc("/max/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/max/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		response, err := uc.GetMaxPrice(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// GET /change/{symbol}
	http.HandleFunc("/change/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/change/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		response, err := uc.GetChangePercent(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// GET /prices?symbols=BTC,ETH
	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "symbols parameter required (e.g., ?symbols=BTC,ETH)", http.StatusBadRequest)
			return
		}

		symbols := strings.Split(symbolsParam, ",")
		for i := range symbols {
			symbols[i] = strings.ToUpper(strings.TrimSpace(symbols[i]))
		}

		prices, err := uc.GetPrices(symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := make(map[string]interface{})
		for symbol, price := range prices {
			response[symbol] = map[string]interface{}{
				"price":        price.Price,
				"last_updated": price.CreatedAt,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// ========== ЗАПУСК СЕРВЕРА ==========

	port := ":8080"
	fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
	fmt.Println("📊 Доступные эндпоинты:")
	fmt.Println("   GET /price/{symbol}        - текущая цена одной валюты")
	fmt.Println("   GET /min/{symbol}          - минимальная цена за всё время")
	fmt.Println("   GET /max/{symbol}          - максимальная цена за всё время")
	fmt.Println("   GET /change/{symbol}       - изменение за час")
	fmt.Println("   GET /prices?symbols=BTC,ETH - цены нескольких валют")

	log.Fatal(http.ListenAndServe(port, nil))
}
