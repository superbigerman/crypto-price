package usecase

import (
	"errors"
	"final/internal/entity"
	"fmt"
)

// ========== ИНТЕРФЕЙСЫ (НА СТОРОНЕ ПОТРЕБИТЕЛЯ) ==========

// PriceRepository — интерфейс для работы с хранилищем цен
type PriceRepository interface {
	GetPrice(symbol string) (entity.Price, error)
	SavePrice(price entity.Price) error
	GetMinPrice(symbol string) (entity.MinPriceResponse, error)
	GetMaxPrice(symbol string) (entity.MaxPriceResponse, error)
	GetChangePercent(symbol string) (entity.ChangePercentResponse, error)
	GetPrices(symbols []string) (map[string]entity.Price, error) // цены нескольких валют
}

// ExternalAPI — интерфейс для внешнего API
type ExternalAPI interface {
	GetRealTimePrice(symbol string) (entity.Price, error)
}

// ========== БИЗНЕС-ЛОГИКА ==========

type PriceUseCase struct {
	repo        PriceRepository
	externalAPI ExternalAPI
}

// NewPriceUseCase — конструктор, внедряет зависимости
func NewPriceUseCase(repo PriceRepository, api ExternalAPI) *PriceUseCase {
	return &PriceUseCase{
		repo:        repo,
		externalAPI: api,
	}
}

// ========== РАБОТА С ЦЕНАМИ ==========

// GetPriceFromDB — получить цену ТОЛЬКО из базы данных
func (uc *PriceUseCase) GetPriceFromDB(symbol string) (entity.Price, error) {
	price, err := uc.repo.GetPrice(symbol)
	if err != nil {
		return entity.Price{}, fmt.Errorf("failed to get price from DB: %w", err)
	}
	return price, nil
}

// FetchAndSave — получить цену из внешнего API и сохранить в БД
func (uc *PriceUseCase) FetchAndSave(symbol string) (entity.Price, error) {
	newPrice, err := uc.externalAPI.GetRealTimePrice(symbol)
	if err != nil {
		return entity.Price{}, fmt.Errorf("failed to fetch from API: %w", err)
	}

	// Сохраняем в БД
	if err := uc.repo.SavePrice(newPrice); err != nil {
		return entity.Price{}, fmt.Errorf("failed to save price to DB: %w", err)
	}

	return newPrice, nil
}

// GetPrice — умный метод: сначала БД, если нет — внешний API
func (uc *PriceUseCase) GetPrice(symbol string) (entity.Price, error) {
	if symbol == "" {
		return entity.Price{}, errors.New("symbol cannot be empty")
	}

	// Приоритет — свои данные (быстрее и надёжнее)
	price, err := uc.GetPriceFromDB(symbol)
	if err == nil {
		return price, nil
	}
	// Fallback — внешний API
	return uc.FetchAndSave(symbol)
}

// GetPrices — получить цены для нескольких валют
func (uc *PriceUseCase) GetPrices(symbols []string) (map[string]entity.Price, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols list cannot be empty")
	}
	return uc.repo.GetPrices(symbols)
}

// SavePrice — сохранить цену (ручное сохранение)
func (uc *PriceUseCase) SavePrice(symbol string, priceValue float64) error {
	if symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	if priceValue < 0 {
		return errors.New("price cannot be negative")
	}

	price, err := entity.NewPrice(symbol, priceValue)
	if err != nil {
		return err
	}

	return uc.repo.SavePrice(*price)
}

// ========== СТАТИСТИКА ==========

// GetMinPrice — минимальная цена за всё время
func (uc *PriceUseCase) GetMinPrice(symbol string) (entity.MinPriceResponse, error) {
	if symbol == "" {
		return entity.MinPriceResponse{}, errors.New("symbol cannot be empty")
	}

	minPrice, err := uc.repo.GetMinPrice(symbol)
	if err != nil {
		return entity.MinPriceResponse{}, fmt.Errorf("failed to get min price: %w", err)
	}

	return minPrice, nil
}

// GetMaxPrice — максимальная цена за всё время
func (uc *PriceUseCase) GetMaxPrice(symbol string) (entity.MaxPriceResponse, error) {
	if symbol == "" {
		return entity.MaxPriceResponse{}, errors.New("symbol cannot be empty")
	}

	maxPrice, err := uc.repo.GetMaxPrice(symbol)
	if err != nil {
		return entity.MaxPriceResponse{}, fmt.Errorf("failed to get max price: %w", err)
	}

	return maxPrice, nil
}

// GetChangePercent — процент изменения цены за час
func (uc *PriceUseCase) GetChangePercent(symbol string) (entity.ChangePercentResponse, error) {
	if symbol == "" {
		return entity.ChangePercentResponse{}, errors.New("symbol cannot be empty")
	}

	changePercent, err := uc.repo.GetChangePercent(symbol)
	if err != nil {
		return entity.ChangePercentResponse{}, fmt.Errorf("failed to get change percent: %w", err)
	}

	return changePercent, nil
}
