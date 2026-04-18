package usecase

import (
	"errors"
	"final/internal/entity"
	"fmt"
	"log"
	"time"
)

type PriceUseCase struct {
	repo        entity.PriceRepository
	externalAPI entity.ExternalAPI
}

func NewPriceUseCase(repo entity.PriceRepository, api entity.ExternalAPI) *PriceUseCase {
	return &PriceUseCase{
		repo:        repo,
		externalAPI: api,
	}
}

// Возвращаем из базы данных
func (uc *PriceUseCase) GetPriceFromDB(symbol string) (entity.Price, error) {
	price, err := uc.repo.GetPrice(symbol)
	if err != nil {
		return entity.Price{}, fmt.Errorf("failed to get price from DB: %w", err)
	}
	return price, nil
}

// Получить цену из внешнего API и сохранить в БД
func (uc *PriceUseCase) FetchAndSave(symbol string) (entity.Price, error) {
	realPrice, err := uc.externalAPI.GetRealTimePrice(symbol)
	if err != nil {
		return entity.Price{}, fmt.Errorf("failed to fetch from API: %w", err)
	}

	newPrice := entity.Price{
		Symbol:    symbol,
		Price:     realPrice,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.SavePrice(newPrice); err != nil {
		log.Printf("WARNING: failed to save price to DB: %v", err)
	}

	return newPrice, nil
}

// Сначала база данных, а потом API
func (uc *PriceUseCase) GetPrice(symbol string) (entity.Price, error) {
	price, err := uc.GetPriceFromDB(symbol)
	if err == nil {
		return price, nil
	}
	return uc.FetchAndSave(symbol)
}

// SavePrice - сохранить цену
func (uc *PriceUseCase) SavePrice(symbol string, priceValue float64) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	if priceValue < 0 {
		return errors.New("price cannot be negative")
	}

	price := entity.Price{
		Symbol:    symbol,
		Price:     priceValue,
		CreatedAt: time.Now(),
	}
	return uc.repo.SavePrice(price)
}

// GetMinPrice - получить минимальную цену
func (uc *PriceUseCase) GetMinPrice(symbol string) (float64, error) {
	if symbol == "" {
		return 0, errors.New("symbol cannot be empty")
	}
	minPrice, err := uc.repo.GetMinPrice(symbol)
	if err != nil {
		return 0, fmt.Errorf("faild to get min price: %w", err)
	}
	return minPrice, nil

}

// GetMaxPrice - получить максимальную цену
func (uc *PriceUseCase) GetMaxPrice(symbol string) (float64, error) {
	if symbol == "" {
		return 0, errors.New("symbol cannot be empty")
	}
	maxPrice, err := uc.repo.GetMaxPrice(symbol)
	if err != nil {
		return 0, fmt.Errorf("faild to get max price: %w", err)
	}
	return maxPrice, nil
}

// GetChangePercent - получить процент изменения за час
// GetChangePercent возвращает процент изменения цены за час
func (uc *PriceUseCase) GetChangePercent(symbol string) (float64, error) {
	if symbol == "" {
		return 0, errors.New("symbol cannot be empty")
	}

	current, err := uc.repo.GetPrice(symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get current price: %w", err)
	}

	hourAgo := time.Now().Add(-1 * time.Hour)
	oldPrice, err := uc.repo.GetPriceAtTime(symbol, hourAgo)
	if err != nil {
		return 0, fmt.Errorf("no price data for 1 hour ago: %w", err)
	}

	if oldPrice.Price == 0 {
		return 0, errors.New("cannot calculate change: price hour ago was zero")
	}

	return ((current.Price - oldPrice.Price) / oldPrice.Price) * 100, nil
}
