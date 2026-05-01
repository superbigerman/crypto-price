package entity

import (
	"errors"
	"time"
)

// Price — основная сущность цены
type Price struct {
	Symbol    string    
	Price     float64   
	CreatedAt time.Time 
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

