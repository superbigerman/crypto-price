package usecase

import (
	"errors"
	"final/internal/entity"
	"fmt"
	"log"
)

type ChangePercentResult struct {
	Symbol       string
	ChagePercent float64
	Diriction    string // "up", "down", "stable"
}

// ========== ИНТЕРФЕЙСЫ ==========

type PriceRepository interface {
	GetPricesLast(symbols []string) ([]entity.Price, error)
	GetMinPrices(symbols []string) ([]entity.Price, error)
	GetMaxPrices(symbols []string) ([]entity.Price, error)
	GetChangePercent(symbols []string) ([]ChangePercentResult, error)
	SavePrice(price entity.Price) error
}

type ExternalAPI interface {
	GetRealTimePrices(symbols []string) ([]entity.Price, error)
}

// ========== БИЗНЕС-ЛОГИКА ==========

type PriceUseCase struct {
	repo        PriceRepository
	externalAPI ExternalAPI
}

func NewPriceUseCase(repo PriceRepository, api ExternalAPI) *PriceUseCase {
	return &PriceUseCase{
		repo:        repo,
		externalAPI: api,
	}
}

// GetPricesLast — возвращает последние цены
func (uc *PriceUseCase) GetPricesLast(symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols list cannot be empty")
	}

	for _, s := range symbols {
		if s == "" {
			return nil, errors.New("symbol cannot be empty")
		}
	}

	dbPrices, err := uc.repo.GetPricesLast(symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices from DB: %w", err)
	}

	if len(dbPrices) == len(symbols) {
		return dbPrices, nil
	}

	apiPrices, err := uc.externalAPI.GetRealTimePrices(symbols)
	if err != nil {
		if len(dbPrices) > 0 {
			return dbPrices, nil
		}
		return nil, fmt.Errorf("failed to get prices from API: %w", err)
	}

	for _, price := range apiPrices {
		if err := uc.repo.SavePrice(price); err != nil {
			log.Printf("WARNING: failed to save price for %s: %v", price.Symbol, err)
		}
	}

	return apiPrices, nil
}

// GetMinPrices — возвращает минимальные цены
func (uc *PriceUseCase) GetMinPrices(symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols list cannot be empty")
	}
	for _, s := range symbols {
		if s == "" {
			return nil, errors.New("symbol cannot be empty")
		}
	}
	minPrices, err := uc.repo.GetMinPrices(symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to get min prices from DB: %w", err)
	}
	if len(minPrices) == 0 {
		return nil, fmt.Errorf("no min prices found for requested symbols")
	}
	return minPrices, nil
}

// GetMaxPrices — возвращает максимальные цены
func (uc *PriceUseCase) GetMaxPrices(symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols list cannot be empty")
	}
	for _, s := range symbols {
		if s == "" {
			return nil, errors.New("symbol cannot be empty")
		}
	}
	maxPrices, err := uc.repo.GetMaxPrices(symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to get max prices from DB: %w", err)
	}
	if len(maxPrices) == 0 {
		return nil, errors.New("no max prices found for requested symbols")
	}
	return maxPrices, nil
}

// GetChangePercent — возвращает изменение за час
func (uc *PriceUseCase) GetChangePercent(symbols []string) ([]ChangePercentResult, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols list cannot be empty")
	}
	for _, s := range symbols {
		if s == "" {
			return nil, errors.New("symbol cannot be empty")
		}
	}
	results, err := uc.repo.GetChangePercent(symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to get change percent from DB: %w", err)
	}
	if len(results) == 0 {
		return nil, errors.New("no change percent data found for requested symbols")
	}
	return results, nil
}
