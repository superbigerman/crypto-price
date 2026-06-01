package usecases

import (
	"context"
	"fmt"
	"log"

	entity "final/internal/entities"
)

// ========== КОНСТРУКТОР ==========

type PriceUseCase struct {
	repo        PriceRepository
	externalAPI ExternalAPI
}

// NewPriceUseCase — конструктор
func NewPriceUseCase(repo PriceRepository, api ExternalAPI) (*PriceUseCase, error) {
	if repo == nil {
		return nil, fmt.Errorf("NewPriceUseCase: PriceRepository cannot be nil")
	}
	if api == nil {
		return nil, fmt.Errorf("NewPriceUseCase: ExternalAPI cannot be nil")
	}
	return &PriceUseCase{
		repo:        repo,
		externalAPI: api,
	}, nil
}

// ========== БИЗНЕС-ЛОГИКА ==========

func (uc *PriceUseCase) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetPricesLast: symbols list cannot be empty")
	}

	// 1. Получаем все валюты, которые есть в таблице currencies
	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		log.Printf("GetPricesLast: failed to get existing symbols: %v", err)
		return nil, fmt.Errorf("GetPricesLast: failed to get existing symbols: %w", err)
	}
	log.Printf("GetPricesLast: existing symbols in DB: %v", existingSymbols)

	// 2. Проверяем, все ли запрошенные валюты есть в currencies
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

	// 3. Если все есть — берём цены из БД
	if allExist {
		log.Printf("GetPricesLast: all symbols found in currencies, fetching from DB")
		prices, err := uc.repo.GetPricesLast(ctx, symbols)
		if err != nil {
			log.Printf("GetPricesLast: DB query failed: %v", err)
			return nil, fmt.Errorf("GetPricesLast: failed to get prices from DB: %w", err)
		}
		log.Printf("GetPricesLast: returning %d prices from DB", len(prices))
		return prices, nil
	}

	// 4. Если нет — идём в API за всеми
	log.Printf("GetPricesLast: fetching from external API for symbols: %v", symbols)
	apiPrices, err := uc.externalAPI.GetRealTimePrices(ctx, symbols)
	if err != nil {
		log.Printf("GetPricesLast: API failed: %v", err)
		return nil, fmt.Errorf("GetPricesLast: failed to get prices from API: %w", err)
	}
	log.Printf("GetPricesLast: API returned %d prices", len(apiPrices))

	// 5. Сохраняем в БД
	if err := uc.repo.SavePrices(ctx, apiPrices); err != nil {
		log.Printf("GetPricesLast: WARNING: failed to save prices to DB: %v", err)
	} else {
		log.Printf("GetPricesLast: successfully saved %d prices to DB", len(apiPrices))
	}

	return apiPrices, nil
}

// GetMinPrices — возвращает минимальные цены
func (uc *PriceUseCase) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetMinPrices: symbols list cannot be empty")
	}

	for _, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetMinPrices: empty symbol")
		}
	}

	existingSymbols, err := uc.repo.GetAllSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetMinPrices: failed to get all symbols: %w", err)
	}

	if len(existingSymbols) == 0 {
		return nil, fmt.Errorf("GetMinPrices: no existing symbols found in database")
	}

	return uc.repo.GetMinPrices(ctx, existingSymbols)
}

// GetMaxPrices — возвращает максимальные цены
func (uc *PriceUseCase) GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetMaxPrices: symbols list cannot be empty")
	}

	for _, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetMaxPrices: empty symbol")
		}
	}

	existingSymbols, err := uc.repo.GetAllSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetMaxPrices: failed to get all symbols: %w", err)
	}

	if len(existingSymbols) == 0 {
		return nil, fmt.Errorf("GetMaxPrices: no existing symbols found in database")
	}

	return uc.repo.GetMaxPrices(ctx, existingSymbols)
}

// GetChangePercent — возвращает изменение за час
func (uc *PriceUseCase) GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetChangePercent: symbols list cannot be empty")
	}

	for _, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetChangePercent: empty symbol")
		}
	}

	existingSymbols, err := uc.repo.GetAllSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetChangePercent: failed to get all symbols: %w", err)
	}

	if len(existingSymbols) == 0 {
		return nil, fmt.Errorf("GetChangePercent: no existing symbols found in database")
	}

	return uc.repo.GetChangePercent(ctx, existingSymbols)
}
