package ports

import (
	"context"

	entity "final/internal/entities"
)

// Client — driven port для внешнего API цен.
type Client interface {
	GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error)
}
