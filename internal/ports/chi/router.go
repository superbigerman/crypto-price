package chi

import (
	"encoding/json"
	"errors"
	entity "final/internal/entities"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type ChiRouter struct {
	mux     *chi.Mux
	useCase PriceUseCase
}

func NewChiRouter(uc PriceUseCase) (*ChiRouter, error) {
	if uc == nil {
		return nil, fmt.Errorf("NewChiRouter: useCase is required")
	}

	rt := &ChiRouter{
		mux:     chi.NewRouter(),
		useCase: uc,
	}
	rt.registerRoutes()
	return rt, nil
}
func (rt *ChiRouter) registerRoutes() {
	rt.mux.Get("/get/prices/last", rt.GetLastPrices)
	rt.mux.Get("/get/prices/min", rt.GetMinPrices)
	rt.mux.Get("/get/prices/max", rt.GetMaxPrices)
	rt.mux.Get("/get/prices/percent", rt.GetChangePrices)
}
func (rt *ChiRouter) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	rt.mux.ServeHTTP(wr, req)
}
func RunServer(uc PriceUseCase) {
	router, err := NewChiRouter(uc)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// ========== ОБЩИЕ ФУНКЦИИ ВАЛИДАЦИИ ==========

// parseAndValidateSymbols парсит и валидирует symbols из запроса
var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)

func parseAndValidateSymbols(req *http.Request) ([]string, error) {
	symbolsParam := req.URL.Query().Get("symbols")
	if symbolsParam == "" {
		return nil, ErrBadRequest
	}

	if strings.Trim(symbolsParam, " ,") == "" {
		return nil, ErrBadRequest
	}

	symbols := strings.Split(symbolsParam, ",")
	validSymbols := make([]string, 0, len(symbols))

	for _, symbol := range symbols {
		symbol = strings.TrimSpace(strings.ToUpper(symbol))
		if symbol == "" {
			continue
		}
		validSymbols = append(validSymbols, symbol) // ← добавить
	}

	if len(validSymbols) == 0 {
		return nil, ErrNotFound
	}

	return validSymbols, nil
}

// ========== DTO ==========

type PriceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}

// ========== КОНВЕРТЕР ==========

func toPriceResponse(prices []entity.Price) []PriceResponse {
	result := make([]PriceResponse, 0, len(prices))
	for _, p := range prices {
		result = append(result, PriceResponse{
			Symbol: p.Symbol,
			Price:  p.Price,
			Time:   p.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}

// ========== GET /prices ==========

func (rt *ChiRouter) GetLastPrices(wr http.ResponseWriter, req *http.Request) {
	validSymbols, err := parseAndValidateSymbols(req)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	prices, err := rt.useCase.GetPricesLast(req.Context(), validSymbols)
	if err != nil {
		http.Error(wr, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(wr).Encode(toPriceResponse(prices))
}

// ========== GET /get/prices/min ==========

func (rt *ChiRouter) GetMinPrices(wr http.ResponseWriter, req *http.Request) {
	validSymbols, err := parseAndValidateSymbols(req)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	prices, err := rt.useCase.GetMinPrices(req.Context(), validSymbols)
	if err != nil {
		http.Error(wr, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(wr).Encode(toPriceResponse(prices))
}

// ========== GET /get/prices/max ==========

func (rt *ChiRouter) GetMaxPrices(wr http.ResponseWriter, req *http.Request) {
	validSymbols, err := parseAndValidateSymbols(req)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	prices, err := rt.useCase.GetMaxPrices(req.Context(), validSymbols)
	if err != nil {
		http.Error(wr, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(wr).Encode(toPriceResponse(prices))
}

// ========== GET /get/prices/percent ==========

func (rt *ChiRouter) GetChangePrices(wr http.ResponseWriter, req *http.Request) {
	validSymbols, err := parseAndValidateSymbols(req)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	changes, err := rt.useCase.GetChangePercent(req.Context(), validSymbols)
	if err != nil {
		http.Error(wr, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(wr).Encode(toPriceResponse(changes))
}
