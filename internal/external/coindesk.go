package external

import (
	"context"
	"encoding/json"
	"final/internal/entity"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type coinDeskResponse struct {
	USD float64 `json:"USD"`
}

type CoinDeskClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewCoinDeskClient() *CoinDeskClient {
	return &CoinDeskClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://min-api.cryptocompare.com",
	}
}

func (c *CoinDeskClient) GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	var result []entity.Price

	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			return nil, ctx.Err() // ← исправлено
		default:
		}

		url := fmt.Sprintf("%s/data/price?fsym=%s&tsyms=USD&relaxedValidation=true",
			c.baseURL,
			strings.ToUpper(symbol))

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			log.Printf("WARN: failed to fetch price for %s: %v", symbol, err)
			continue
		}

		if resp.StatusCode >= 500 {
			resp.Body.Close()
			return nil, fmt.Errorf("API server error for %s: %d", symbol, resp.StatusCode)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("WARN: API returned %d for %s, body: %s", resp.StatusCode, symbol, string(body))
			resp.Body.Close()
			continue
		}

		var apiResponse coinDeskResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			log.Printf("WARN: failed to decode response for %s: %v", symbol, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		usdPrice := apiResponse.USD
		if usdPrice == 0 {
			// цена 0 допустима, не пропускаем — но ты решила оставить проверку
			// можешь убрать эту проверку, если хочешь
			continue
		}

		newPrice, err := entity.NewPrice(strings.ToUpper(symbol), usdPrice)
		if err != nil {
			continue
		}

		result = append(result, *newPrice)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no prices found for any of the requested symbols")
	}

	return result, nil
}
