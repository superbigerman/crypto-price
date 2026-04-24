package entity

import (
	"errors"
	"time"
)

// Price — основная сущность цены
type Price struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// NewPrice — конструктор для Price
func NewPrice(symbol string, price float64) (*Price, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	return &Price{
		Symbol:    symbol,
		Price:     price,
		CreatedAt: time.Now(),
	}, nil
}

// MinPriceResponse — ответ для минимальной цены
type MinPriceResponse struct {
	Symbol     string    `json:"symbol"`
	MinPrice   float64   `json:"min_price"`
	AchievedAt time.Time `json:"achieved_at"`
}

// NewMinPriceResponse — конструктор для MinPriceResponse
func NewMinPriceResponse(symbol string, minPrice float64, achievedAt time.Time) *MinPriceResponse {
	return &MinPriceResponse{
		Symbol:     symbol,
		MinPrice:   minPrice,
		AchievedAt: achievedAt,
	}
}

// MaxPriceResponse — ответ для максимальной цены
type MaxPriceResponse struct {
	Symbol     string    `json:"symbol"`
	MaxPrice   float64   `json:"max_price"`
	AchievedAt time.Time `json:"achieved_at"`
}

// NewMaxPriceResponse — конструктор для MaxPriceResponse
func NewMaxPriceResponse(symbol string, maxPrice float64, achievedAt time.Time) *MaxPriceResponse {
	return &MaxPriceResponse{
		Symbol:     symbol,
		MaxPrice:   maxPrice,
		AchievedAt: achievedAt,
	}
}

// ChangePercentResponse — ответ для изменения за час
type ChangePercentResponse struct {
	Symbol        string    `json:"symbol"`
	ChangePercent float64   `json:"change_percent"`
	Direction     string    `json:"direction"`
	Period        string    `json:"period"`
	CalculatedAt  time.Time `json:"calculated_at"`
}

// NewChangePercentResponse — конструктор для ChangePercentResponse
func NewChangePercentResponse(symbol string, changePercent float64, direction string) *ChangePercentResponse {
	return &ChangePercentResponse{
		Symbol:        symbol,
		ChangePercent: changePercent,
		Direction:     direction,
		Period:        "1h",
		CalculatedAt:  time.Now(),
	}
}
