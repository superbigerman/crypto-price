package ports

import (
	"context"

	entity "final/internal/entities"
)

// PriceUseCase — driving port для сценариев работы с ценами.
type PriceUseCase interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
}
