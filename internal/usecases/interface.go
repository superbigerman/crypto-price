package usecases

import (
	"context"
	entity "final/internal/entities"
)

// PriceRepository
type PriceRepository interface {
	SavePrices(ctx context.Context, prices []entity.Price) error
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
	AddCurrencies(ctx context.Context, symbols []string) error
	GetAllSymbols(ctx context.Context) ([]string, error)
	GetExistingSymbols(ctx context.Context, symbols []string) ([]string, error)
}

// PriceUseCase
type PriceUseCase interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
}

// PriceClient
type PriceClient interface {
	GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error)
}
