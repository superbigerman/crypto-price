package usecases

import (
	"context"
	"fmt"
	"log"

	entity "final/internal/entities"
)

type PriceUseCaseImpl struct {
	repo   PriceRepository
	client PriceClient
}

func NewPriceUseCase(repo PriceRepository, client PriceClient) (PriceUseCase, error) {
	if repo == nil {
		return nil, fmt.Errorf("NewPriceUseCase: PriceRepository cannot be nil")
	}
	if client == nil {
		return nil, fmt.Errorf("NewPriceUseCase: PriceClient cannot be nil")
	}
	return &PriceUseCaseImpl{
		repo:   repo,
		client: client,
	}, nil
}

func (uc *PriceUseCaseImpl) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
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

	// 4. Сохраняем цены
	if err := uc.repo.SavePrices(ctx, apiPrices); err != nil {
		log.Printf("WARNING: failed to save prices: %v", err)
	}

	return apiPrices, nil
}

func (uc *PriceUseCaseImpl) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
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

func (uc *PriceUseCaseImpl) GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
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

func (uc *PriceUseCaseImpl) GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error) {
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
