package output

import (
	"context"

	entity "final/internal/entities"
)

// PriceProvider — выходной порт (driven): внешний источник актуальных цен.
type PriceProvider interface {
	GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error)
}
