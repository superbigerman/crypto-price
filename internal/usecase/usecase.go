package usecase

import (
	"errors"
	"final/internal/entity"
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

// GetPrice - получить цену (любой валюты!)
func (uc *PriceUseCase) GetPrice(symbol string) (entity.Price, error) {
	// 1. Сначала пробуем получить из базы
	price, err := uc.repo.GetPrice(symbol)
	if err == nil {
		return price, nil
	}

	// 2. В базе нет — идём во внешний API
	realPrice, err := uc.externalAPI.GetRealTimePrice(symbol)
	if err != nil {
		return entity.Price{}, errors.New("price for " + symbol + " not available")
	}

	// 3. Сохраняем в базу для следующих запросов
	newPrice := entity.Price{
		Symbol:    symbol,
		Price:     realPrice,
		CreatedAt: time.Now(),
	}
	uc.repo.SavePrice(newPrice)

	return newPrice, nil
}

// SavePrice - сохранить цену
func (uc *PriceUseCase) SavePrice(symbol string, priceValue float64) error {
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
	if symbol != "BTC" && symbol != "ETH" {
		return 0, errors.New("min price available only for BTC and ETH")
	}

	prices, err := uc.repo.GetAllPrices(symbol)
	if err != nil {
		return 0, err
	}

	if len(prices) == 0 {
		return 0, errors.New("no prices found for " + symbol)
	}

	minPrice := prices[0].Price
	for _, p := range prices {
		if p.Price < minPrice {
			minPrice = p.Price
		}
	}

	return minPrice, nil
}

// GetMaxPrice - получить максимальную цену
func (uc *PriceUseCase) GetMaxPrice(symbol string) (float64, error) {
	if symbol != "BTC" && symbol != "ETH" {
		return 0, errors.New("max price available only for BTC and ETH")
	}

	prices, err := uc.repo.GetAllPrices(symbol)
	if err != nil {
		return 0, err
	}

	if len(prices) == 0 {
		return 0, errors.New("no prices found for " + symbol)
	}

	maxPrice := prices[0].Price
	for _, p := range prices {
		if p.Price > maxPrice {
			maxPrice = p.Price
		}
	}

	return maxPrice, nil
}

// GetChangePercent - получить процент изменения за час
func (uc *PriceUseCase) GetChangePercent(symbol string) (float64, error) {
	if symbol != "BTC" && symbol != "ETH" {
		return 0, errors.New("change percent available only for BTC and ETH")
	}

	prices, err := uc.repo.GetAllPrices(symbol)
	if err != nil {
		return 0, err
	}

	if len(prices) == 0 {
		return 0, errors.New("no prices found for " + symbol)
	}

	currentPrice := prices[len(prices)-1].Price

	hourAgo := time.Now().Add(-1 * time.Hour)
	var hourAgoPrice float64
	found := false

	for i := len(prices) - 1; i >= 0; i-- {
		if prices[i].CreatedAt.Before(hourAgo) || prices[i].CreatedAt.Equal(hourAgo) {
			hourAgoPrice = prices[i].Price
			found = true
			break
		}
	}

	if !found {
		hourAgoPrice = prices[0].Price
	}

	if hourAgoPrice == 0 {
		return 0, errors.New("cannot calculate change")
	}

	return ((currentPrice - hourAgoPrice) / hourAgoPrice) * 100, nil
}
