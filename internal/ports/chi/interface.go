package chi

import (
	"context"
	entity "final/internal/entities"
)

// ========== ИНТЕРФЕЙС ==========

type PriceUseCase interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
}
