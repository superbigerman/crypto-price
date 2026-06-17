package coindesk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	entity "final/internal/entities"
)

type CoinDeskClient struct {
	httpClient *http.Client
	baseURL    string
	relaxed    bool
	tsyms      string
	apiKey     string
}

func NewCoinDeskClient(baseURL string, timeout time.Duration, relaxed bool, tsyms string, apiKey string) *CoinDeskClient {
	return &CoinDeskClient{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		relaxed:    relaxed,
		tsyms:      tsyms,
		apiKey:     apiKey,
	}
}

func (c *CoinDeskClient) GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetRealTimePrices: symbols list cannot be empty")
	}

	// Строим URL
	rawURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("GetRealTimePrices: invalid base URL: %w", err)
	}

	rawURL.Path = "/data/pricemulti"

	query := rawURL.Query()
	query.Set("fsyms", strings.Join(symbols, ","))
	query.Set("tsyms", c.tsyms)
	query.Set("relaxedValidation", fmt.Sprintf("%v", c.relaxed))
	rawURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetRealTimePrices: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Apikey "+c.apiKey)

	// ========== ВЫПОЛНЯЕМ ЗАПРОС С ПОВТОРОМ ПРИ 429 ==========
	var resp *http.Response
	maxRetries := 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			log.Printf("GetRealTimePrices: network error (attempt %d): %v", attempt+1, err)
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("GetRealTimePrices: network error after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		// 429 Too Many Requests — повторяем
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("GetRealTimePrices: rate limit (429), retrying (attempt %d)", attempt+1)
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("GetRealTimePrices: rate limit exceeded after %d attempts", maxRetries)
			}
			time.Sleep(time.Duration(attempt+2) * time.Second)
			continue
		}

		// 500+ — ошибка сервера
		if resp.StatusCode >= 500 {
			log.Printf("GetRealTimePrices: API server error: %d", resp.StatusCode)
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("GetRealTimePrices: API server error: %d", resp.StatusCode)
			}
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		// 4xx — клиентская ошибка (не повторяем)
		if resp.StatusCode != http.StatusOK {
			log.Printf("GetRealTimePrices: API returned status %d for symbols %v", resp.StatusCode, symbols)
			return nil, fmt.Errorf("GetRealTimePrices: API returned status %d", resp.StatusCode)
		}

		// Успех — выходим из цикла
		break
	}

	// ========== ПАРСИНГ JSON ==========
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("GetRealTimePrices: failed to decode response: %v", err)
		return nil, fmt.Errorf("GetRealTimePrices: failed to decode response: %w", err)
	}

	var prices []entity.Price
	for symbol, val := range data {
		usd, ok := val["USD"]
		if !ok {
			log.Printf("GetRealTimePrices: no USD price for %s", symbol)
			continue
		}
		p, err := entity.NewPrice(symbol, usd)
		if err != nil {
			log.Printf("GetRealTimePrices: failed to create price for %s: %v", symbol, err)
			continue
		}
		prices = append(prices, *p)
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("GetRealTimePrices: no prices found for %v", symbols)
	}

	log.Printf("GetRealTimePrices: successfully fetched %d prices", len(prices))
	return prices, nil
}
