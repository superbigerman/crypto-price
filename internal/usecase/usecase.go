package usecase

import (
	"context"
	"final/internal/entity"
	"fmt"
	"log"
)

// ========== ИНТЕРФЕЙСЫ ==========

type PriceRepository interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]float64, error)
	SavePrices(ctx context.Context, prices []entity.Price) error
	GetExistingSymbols(ctx context.Context, symbols []string) ([]string, error)
}

type ExternalAPI interface {
	GetRealTimePrices(ctx context.Context, symbols []string) ([]entity.Price, error)
}

// ========== БИЗНЕС-ЛОГИКА ==========

type PriceUseCase struct {
	repo        PriceRepository
	externalAPI ExternalAPI
}

// GetPricesLast — возвращает последние цены
func (uc *PriceUseCase) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		return nil, err
	}

	// Если все валюты уже есть в таблице currencies → берём из БД
	if len(existingSymbols) == len(symbols) {
		return uc.repo.GetPricesLast(ctx, symbols)
	}

	// Если какой‑то валюты нет → идём во внешний API за всеми
	apiPrices, err := uc.externalAPI.GetRealTimePrices(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices from API: %w", err)
	}

	// Сохраняем цены (репозиторий сам разберётся с добавлением валют)
	if err := uc.repo.SavePrices(ctx, apiPrices); err != nil {
		log.Printf("WARNING: failed to save prices: %v", err)
	}

	return apiPrices, nil
}

// GetMinPrices — возвращает минимальные цены
func (uc *PriceUseCase) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		return nil, err
	}
	if len(existingSymbols) == 0 {
		return nil, fmt.Errorf("no existingSymbols")
	}
	return uc.repo.GetMinPrices(ctx, existingSymbols)
}

// GetMaxPrices — возвращает максимальные цены
func (uc *PriceUseCase) GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		return nil, err
	}
	if len(existingSymbols) == 0 {
		return nil, fmt.Errorf("no existingSymbols")
	}
	return uc.repo.GetMaxPrices(ctx, existingSymbols)
}

// GetChangePercent — возвращает изменение за час
func (uc *PriceUseCase) GetChangePercent(ctx context.Context, symbols []string) ([]float64, error) {
	existingSymbols, err := uc.repo.GetExistingSymbols(ctx, symbols)
	if err != nil {
		return nil, err
	}
	if len(existingSymbols) == 0 {
		return nil, fmt.Errorf("no existing symbols provided")
	}
	return uc.repo.GetChangePercent(ctx, existingSymbols)
}
