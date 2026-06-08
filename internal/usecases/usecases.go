package usecases

import (
	"context"
	"fmt"
	"log"

	entity "final/internal/entities"
	"final/internal/ports/input"
	"final/internal/ports/output"
)

var _ input.PriceService = (*PriceUseCase)(nil)

// ========== КОНСТРУКТОР ==========

type PriceUseCase struct {
	repo     output.PriceRepository
	provider output.PriceClient
}

// NewPriceUseCase — конструктор
func NewPriceUseCase(repo output.PriceRepository, provider output.PriceProvider) (input.PriceService, error) {
	if repo == nil {
		return nil, fmt.Errorf("NewPriceUseCase: PriceRepository cannot be nil")
	}
	if provider == nil {
		return nil, fmt.Errorf("NewPriceUseCase: PriceProvider cannot be nil")
	}
	return &PriceUseCase{
		repo:     repo,
		provider: provider,
	}, nil
}

// ========== БИЗНЕС-ЛОГИКА ==========

func (uc *PriceUseCase) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
	// 1. Получаем существующие валюты
	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		return nil, err
	}

	// 2. Если все есть — берём из БД
	if len(existingSymbols) == len(symbols) {
		return uc.repo.GetPricesLast(ctx, symbols)
	}

	// 3. Идём в API
	apiPrices, err := uc.client.GetRealTimePrices(ctx, symbols)
	if err != nil {
		return nil, err
	}

	// 4. Сохраняем новые ВАЛЮТЫ (отдельно!)
	newSymbols := extractSymbols(apiPrices) // ["XRP", "SOL"...]
	if err := uc.repo.AddCurrencies(ctx, newSymbols); err != nil {
		log.Printf("WARNING: failed to add currencies: %v", err)
	}

	// 5. Сохраняем ЦЕНЫ (отдельно!)
	if err := uc.repo.SavePrices(ctx, apiPrices); err != nil {
		log.Printf("WARNING: failed to save prices: %v", err)
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

// extractSymbols извлекает символы валют из слайса цен
func extractSymbols(prices []entity.Price) []string {
	symbols := make([]string, len(prices))
	for i, p := range prices {
		symbols[i] = p.Symbol
	}
	return symbols
}
