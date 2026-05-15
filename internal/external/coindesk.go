package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"final/internal/config"
	"final/internal/entity"
)

type CoinDeskClient struct {
	httpClient *http.Client
	baseURL    string
	relaxed    bool
	tsyms      string
}

func NewCoinDeskClient(cfg *config.Config) *CoinDeskClient {
	if cfg == nil {
		cfg = config.Load()
	}
	return &CoinDeskClient{
		httpClient: &http.Client{
			Timeout: cfg.ExternalAPITimeout,
		},
		baseURL: cfg.ExternalAPIBaseURL,
		relaxed: cfg.ExternalAPIRelaxed,
		tsyms:   cfg.ExternalAPITSyms,
	}
}

func (c *CoinDeskClient) GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	fsyms := strings.Join(symbols, ",")
	url := fmt.Sprintf("%s/data/pricemulti?fsyms=%s&tsyms=%s&relaxedValidation=%v",
		c.baseURL, fsyms, c.tsyms, c.relaxed)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var prices []entity.Price
	for symbol, val := range data {
		usd, ok := val["USD"]
		if !ok {
			continue
		}
		p, err := entity.NewPrice(symbol, usd)
		if err != nil {
			continue
		}
		prices = append(prices, *p)
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("no prices found")
	}

	return prices, nil
}
