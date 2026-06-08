package output

import (
	"context"
	entity "final/internal/entities"
)

type PriceRepository interface {
	// Работа с таблицей prices
	SavePrices(ctx context.Context, prices []entity.Price) error
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)

	// Работа с таблицей currencies (ОТДЕЛЬНО!)
	AddCurrencies(ctx context.Context, symbols []string) error
	GetAllSymbols(ctx context.Context) ([]string, error)
	GetExistingSymbols(ctx context.Context, symbols []string) ([]string, error)
}
