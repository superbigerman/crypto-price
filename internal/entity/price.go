package entity

import "time"

type Price struct {
	Symbol    string
	Price     float64
	CreatedAt time.Time
}

type MinPriceResponse struct {
	Symbol   string
	MinPrice float64
	UpdateAt time.Time
}

type MaxPriceRespose struct {
	Symbol   string
	MaxPrice float64
	UpdateAt time.Time
}

type ChangeResponse struct {
	Symbol        string
	Changepercent float64
	Direction     string
	Period        string
	UpdateAt      time.Time
}

type PriceRepository interface {
	GetPrice(symbol string) (Price, error)
	SavePrice(price Price) error
	GetAllPrices(symbol string) ([]Price, error)
}

type ExternalAPI interface {
	GetRealTimePrice(symbol string) (float64, error)
}

// TODO конструктор , ошибки из каждой фунции должна возращать ошибку,
