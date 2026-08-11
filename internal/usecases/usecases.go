package usecases

import (
	"context"
	"fmt"
	"log"

	entity "final/internal/entities"
)

// ========== СТРУКТУРА ==========

type PriceUseCaseImpl struct {
	repo        PriceRepository
	externalAPI PriceClient
}

// ========== КОНСТРУКТОР ==========

func NewPriceUseCase(repo PriceRepository, api PriceClient) (*PriceUseCaseImpl, error) {
	if repo == nil {
		return nil, fmt.Errorf("NewPriceUseCase: PriceRepository cannot be nil")
	}
	if api == nil {
		return nil, fmt.Errorf("NewPriceUseCase: ExternalAPI cannot be nil")
	}

	return &PriceUseCaseImpl{
		repo:        repo,
		externalAPI: api,
	}, nil
}

// ========== БИЗНЕС-ЛОГИКА ==========

func (uc *PriceUseCaseImpl) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetPricesLast: symbols list cannot be empty")
	}

	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		log.Printf("GetPricesLast: failed to get existing symbols: %v", err)
		return nil, fmt.Errorf("GetPricesLast: failed to get existing symbols: %w", err)
	}

	log.Printf("GetPricesLast: existing symbols in DB: %v", existingSymbols)

	allExist := true
	for _, s := range symbols {
		found := false
		for _, e := range existingSymbols {
			if s == e {
				found = true
				break
			}
		}
		if !found {
			allExist = false
			log.Printf("GetPricesLast: symbol %s not found in currencies", s)
			break
		}
	}

	if allExist {
		log.Printf("GetPricesLast: all symbols found, fetching from DB")
		prices, err := uc.repo.GetPricesLast(ctx, symbols)
		if err != nil {
			log.Printf("GetPricesLast: DB query failed: %v", err)
			return nil, fmt.Errorf("GetPricesLast: failed to get prices from DB: %w", err)
		}
		log.Printf("GetPricesLast: returning %d prices from DB", len(prices))
		return prices, nil
	}

	log.Printf("GetPricesLast: fetching from external API for symbols: %v", symbols)
	apiPrices, err := uc.externalAPI.GetRealTimePrices(ctx, symbols)
	if err != nil {
		log.Printf("GetPricesLast: API failed: %v", err)
		return nil, fmt.Errorf("GetPricesLast: failed to get prices from API: %w", err)
	}
	log.Printf("GetPricesLast: API returned %d prices", len(apiPrices))

	if err := uc.repo.SavePrices(ctx, apiPrices); err != nil {
		log.Printf("GetPricesLast: WARNING: failed to save prices to DB: %v", err)
	} else {
		log.Printf("GetPricesLast: successfully saved %d prices to DB", len(apiPrices))
	}

	return apiPrices, nil
}

func (uc *PriceUseCaseImpl) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetMinPrices: symbols list cannot be empty")
	}
	return uc.repo.GetMinPrices(ctx, symbols)
}

func (uc *PriceUseCaseImpl) GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetMaxPrices: symbols list cannot be empty")
	}
	return uc.repo.GetMaxPrices(ctx, symbols)
}
func (uc *PriceUseCaseImpl) GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetChangePercent: symbols list cannot be empty")
	}

	// Вызываем GetPricesLast, чтобы добавить недостающие валюты
	_, err := uc.GetPricesLast(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("GetChangePercent: failed to ensure symbols exist: %w", err)
	}

	// Теперь все валюты есть в БД
	return uc.repo.GetChangePercent(ctx, symbols)
}

func (uc *PriceUseCaseImpl) UpdateAllPrices(ctx context.Context) error {
	symbols, err := uc.repo.GetAllSymbols(ctx)
	if err != nil {
		return fmt.Errorf("UpdateAllPrices: %w", err)
	}

	if len(symbols) == 0 {
		return nil
	}

	prices, err := uc.externalAPI.GetRealTimePrices(ctx, symbols)
	if err != nil {
		return fmt.Errorf("UpdateAllPrices: %w", err)
	}

	return uc.repo.SavePrices(ctx, prices)
}

func (uc *PriceUseCaseImpl) GetAllSymbols(ctx context.Context) ([]string, error) {
	return uc.repo.GetAllSymbols(ctx)
}
