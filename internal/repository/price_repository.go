package repository

import "final/internal/entity"

type PreceRepository interface {
	GetPrice(symbol string) (float64, error)
	SavePrice(price entity.Price) error
}
