package usecase

import (
	"context"
	"fmt"
	"log"

	entity "final/internal/entities"
)

// ========== ИНТЕРФЕЙСЫ ==========

type PriceRepository interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
	SavePrices(ctx context.Context, prices []entity.Price) error
	GetAllSymbols(ctx context.Context) ([]string, error)
}

type ExternalAPI interface {
	GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error)
}

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

// GetPricesLast — возвращает последние цены
func (uc *PriceUseCase) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetPricesLast: symbols list cannot be empty")
	}

	// Проверяем, есть ли пустые строки
	for i, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetPricesLast: empty symbol at index %d", i)
		}
	}

	// 1. Пробуем получить из БД
	prices, err := uc.repo.GetPricesLast(ctx, symbols)
	if err == nil && len(prices) == len(symbols) {
		return prices, nil
	}
	if err != nil {
		log.Printf("GetPricesLast: DB query failed: %v", err)
	}

	// 2. Если нет — идём в API
	apiPrices, err := uc.externalAPI.GetRealTimePrices(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("GetPricesLast: failed to get prices from API: %w", err)
	}

	// 3. Сохраняем в БД
	if err := uc.repo.SavePrices(ctx, apiPrices); err != nil {
		log.Printf("GetPricesLast: WARNING: failed to save prices to DB: %v", err)
	}

	return apiPrices, nil
}

// GetMinPrices — возвращает минимальные цены
func (uc *PriceUseCase) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("GetMinPrices: symbols list cannot be empty")
	}

	for i, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetMinPrices: empty symbol at index %d", i)
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

	for i, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetMaxPrices: empty symbol at index %d", i)
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

	for i, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("GetChangePercent: empty symbol at index %d", i)
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
