package usecases

import (
	"context"
	entity "final/internal/entities"
)

// PriceRepository — интерфейс для работы с хранилищем цен
type PriceRepository interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
	SavePrices(ctx context.Context, prices []entity.Price) error
	GetAllSymbols(ctx context.Context) ([]string, error)
	GetExistingSymbols(ctx context.Context, symbols []string) ([]string, error)
}

// ExternalAPI — интерфейс для внешнего API
type ExternalAPI interface {
	GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error)
}
