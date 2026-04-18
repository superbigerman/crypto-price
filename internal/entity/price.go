package entity

import "time"

type Price struct {
	Symbol    string
	Price     float64
	CreatedAt time.Time
}

type PriceRepository interface {
	GetPrice(symbol string) (Price, error)
	SavePrice(price Price) error
	GetMinPrice(symbol string) (float64, error)
	GetMaxPrice(symbol string) (float64, error)
	GetPriceAtTime(symbol string, timestamp time.Time) (Price, error)
}

type ExternalAPI interface {
	GetRealTimePrice(symbol string) (float64, error)
}
