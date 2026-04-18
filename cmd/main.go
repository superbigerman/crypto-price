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

func (m *mockRepo) GetMinPrice(symbol string) (float64, error) {
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return 0, fmt.Errorf("no prices for %s", symbol)
	}
	minPrice := prices[0].Price
	for _, p := range prices {
		if p.Price < minPrice {
			minPrice = p.Price
		}
	}
	return minPrice, nil
}

func (m *mockRepo) GetMaxPrice(symbol string) (float64, error) {
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return 0, fmt.Errorf("no prices for %s", symbol)
	}
	maxPrice := prices[0].Price
	for _, p := range prices {
		if p.Price > maxPrice {
			maxPrice = p.Price
		}
	}
	return maxPrice, nil
}

func (m *mockRepo) GetPriceAtTime(symbol string, timestamp time.Time) (entity.Price, error) {
	prices, ok := m.prices[symbol]
	if !ok || len(prices) == 0 {
		return entity.Price{}, fmt.Errorf("no prices for %s", symbol)
	}
	for i := len(prices) - 1; i >= 0; i-- {
		if prices[i].CreatedAt.Before(timestamp) || prices[i].CreatedAt.Equal(timestamp) {
			return prices[i], nil
		}
	}
	return entity.Price{}, fmt.Errorf("no price before %v", timestamp)
}

type mockExternalAPI struct{}

func (m *mockExternalAPI) GetRealTimePrice(symbol string) (float64, error) {
	mockPrices := map[string]float64{
		"BTC":  45000,
		"ETH":  3200,
		"XRP":  0.5,
		"DOGE": 0.08,
		"SOL":  100,
	}
	if price, ok := mockPrices[symbol]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("price for %s not available", symbol)
}

func main() {
	repo := &mockRepo{}
	externalAPI := &mockExternalAPI{}
	uc := usecase.NewPriceUseCase(repo, externalAPI)

	// Тестовые данные
	uc.SavePrice("BTC", 44000)
	time.Sleep(2 * time.Second)
	uc.SavePrice("BTC", 44500)
	time.Sleep(2 * time.Second)
	uc.SavePrice("BTC", 45000)
	time.Sleep(2 * time.Second)
	uc.SavePrice("BTC", 44800)

	uc.SavePrice("ETH", 3100)
	time.Sleep(2 * time.Second)
	uc.SavePrice("ETH", 3200)
	time.Sleep(2 * time.Second)
	uc.SavePrice("ETH", 3250)

	fmt.Println("✅ Тестовые данные загружены")

	// HTTP handlers
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

	http.HandleFunc("/min/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/min/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		minPrice, err := uc.GetMinPrice(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol":    symbol,
			"min_price": minPrice,
		})
	})

	http.HandleFunc("/max/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/max/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		maxPrice, err := uc.GetMaxPrice(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol":    symbol,
			"max_price": maxPrice,
		})
	})

	http.HandleFunc("/change/", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/change/")
		symbol = strings.ToUpper(symbol)

		if symbol == "" {
			http.Error(w, "symbol required", http.StatusBadRequest)
			return
		}

		changePercent, err := uc.GetChangePercent(symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		direction := "stable"
		if changePercent > 0 {
			direction = "up"
		} else if changePercent < 0 {
			direction = "down"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol":         symbol,
			"change_percent": changePercent,
			"direction":      direction,
			"period":         "1h",
		})
	})

	port := ":8080"
	fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
	fmt.Println("📊 Доступные эндпоинты:")
	fmt.Println("   GET /price/{symbol}  - текущая цена")
	fmt.Println("   GET /min/{symbol}    - минимальная цена за всё время")
	fmt.Println("   GET /max/{symbol}    - максимальная цена за всё время")
	fmt.Println("   GET /change/{symbol} - изменение за час")

	log.Fatal(http.ListenAndServe(port, nil))
}
