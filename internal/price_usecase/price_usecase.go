package usecase

import (
	"errors"
	"final/internal/entity"
	"final/internal/repository"
)

type PriceUseCase struct {
	repo repository.PriceRepository
}

func NewPriceCase(repo repository.PriceRepository) *PriceUseCase {
	return &PriceUseCase{repo: repo}
}

func (uc *PriceUseCase) GetPrice(symbol string) (float64, error) {
	if symbol != "BTC" && symbol != "EHT" {
		return 0, errors.New("only BTC or EHT")
	}
	return uc.repo.GetPrice(symbol)
}

func (uc *PriceUseCase) SavePrice(symbol string, price float64) error {
	if symbol != "BTC" && symbol != "EHT" {
		return errors.New("only BTC or EHT")
	}
	if price <= 0 {

		return errors.New("price must be positive")
	}
	return uc.repo.SavePrice(entity.Price{Symbol: symbol, Price: price})
}
