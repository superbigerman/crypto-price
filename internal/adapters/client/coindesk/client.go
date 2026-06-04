package coindesk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"final/config"
	entity "final/internal/entities"
	"final/internal/ports"
)

var _ ports.Client = (*CoinDeskClient)(nil)

type CoinDeskClient struct {
	httpClient *http.Client
	baseURL    string
	relaxed    bool
	tsyms      string
}

func NewCoinDeskClient(cfg *config.Config) (*CoinDeskClient, error) {
	if cfg == nil {
		cfg = config.Load()
	}
	if cfg == nil {
		return nil, fmt.Errorf("NewCoinDeskClient: config is nil after load")
	}
	return &CoinDeskClient{
		httpClient: &http.Client{
			Timeout: cfg.ExternalAPITimeout,
		},
		baseURL: cfg.ExternalAPIBaseURL,
		relaxed: cfg.ExternalAPIRelaxed,
		tsyms:   cfg.ExternalAPITSyms,
	}, nil
}

func (c *CoinDeskClient) GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetRealTimePrices: symbols list cannot be empty")
	}

	rawURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("GetRealTimePrice: invalid base URL: %w", err)
	}
	rawURL.Path = "/datapricemulti"

	query := rawURL.Query()
	query.Set("tsyms", strings.Join(symbols, ","))
	query.Set("tsyms", c.tsyms)
	query.Set("relaxedValidation", fmt.Sprintf("%v", c.relaxed))
	rawURL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetRealTimePrice: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetRealTimePrices: API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetRealTimePrices: API returned status %d", resp.StatusCode)
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("GetRealTimePrices: failed to decode response: %w", err)
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
		return nil, fmt.Errorf("GetRealTimePrices: no prices found for %v", symbols)
	}

	return prices, nil
}
