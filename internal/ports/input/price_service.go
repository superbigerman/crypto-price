package input

import (
	"context"

	entity "final/internal/entities"
)

// PriceService — входной порт (driving): контракт прикладного слоя для внешних адаптеров.
type PriceService interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
}
