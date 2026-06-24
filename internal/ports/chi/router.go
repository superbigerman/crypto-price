package chi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	entity "final/internal/entities"

	"github.com/go-chi/chi/v5"
)

// ========== ИНТЕРФЕЙС ==========

type PriceUseCase interface {
	GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error)
	GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error)
}

// ========== ХЕНДЛЕР ==========

type PriceHandler struct {
	useCase PriceUseCase
}

func NewPriceHandler(uc PriceUseCase) (*PriceHandler, error) {
	if uc == nil {
		return nil, fmt.Errorf("NewPriceHandler: useCase cannot be nil")
	}
	return &PriceHandler{useCase: uc}, nil
}

// ========== ОБЩИЕ ФУНКЦИИ ВАЛИДАЦИИ ==========

// parseAndValidateSymbols парсит и валидирует symbols из запроса
func parseAndValidateSymbols(r *http.Request) ([]string, error) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		return nil, fmt.Errorf("query parameter 'symbols' is required, example: /prices?symbols=BTC,ETH,ADA")
	}

	symbols := strings.Split(symbolsParam, ",")
	if len(symbols) == 0 {
		return nil, fmt.Errorf("empty symbols list")
	}

	validSymbols := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		symbol = strings.TrimSpace(strings.ToUpper(symbol))
		if symbol == "" {
			continue
		}
		if len(symbol) < 2 || len(symbol) > 10 {
			return nil, fmt.Errorf("invalid symbol format: %s", symbol)
		}
		validSymbols = append(validSymbols, symbol)
	}

	if len(validSymbols) == 0 {
		return nil, fmt.Errorf("no valid symbols provided")
	}

	if len(validSymbols) > 50 {
		return nil, fmt.Errorf("too many symbols, max 50 allowed")
	}

	return validSymbols, nil
}

// writeSuccess отправляет успешный JSON-ответ с метаданными
func writeSuccess(w http.ResponseWriter, data interface{}, meta map[string]interface{}) {
	response := map[string]interface{}{
		"data": data,
		"meta": meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: failed to encode response: %v", err)
	}
}

// writeError отправляет JSON-ошибку
func writeError(w http.ResponseWriter, status int, errorType string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errorType,
		"message": message,
	})
}

// handleServiceError обрабатывает ошибки от бизнес-логики
func handleServiceError(w http.ResponseWriter, err error) {
	// Логируем реальную ошибку
	log.Printf("ERROR: %v", err)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusGatewayTimeout, "request timeout",
			"upstream service did not respond in time")
	case strings.Contains(err.Error(), "not found"):
		writeError(w, http.StatusNotFound, "not found", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error",
			"failed to process request")
	}
}

// ========== GET /prices ==========

// ========== GET /prices ==========

func (h *PriceHandler) GetLastPrices(w http.ResponseWriter, r *http.Request) {
	// 1. Проверка метода HTTP
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed",
			"only GET method is supported")
		return
	}

	// 2. Парсинг и валидация symbols
	validSymbols, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	// 3. Таймаут контекста
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 4. Вызов бизнес-логики
	prices, err := h.useCase.GetPricesLast(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// 5. Проверка пустого ответа
	if len(prices) == 0 {
		writeError(w, http.StatusNotFound, "no data",
			fmt.Sprintf("no prices found for symbols: %v", validSymbols))
		return
	}

	// 6. Успешный ответ с метаданными
	writeSuccess(w, prices, map[string]interface{}{
		"count":        len(prices),
		"symbols":      validSymbols,
		"requested_at": time.Now().UTC().Format(time.RFC3339),
	})
}

// ========== GET /min ==========

func (h *PriceHandler) GetMinPrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed",
			"only GET method is supported")
		return
	}

	validSymbols, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	prices, err := h.useCase.GetMinPrices(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if len(prices) == 0 {
		writeError(w, http.StatusNotFound, "no data",
			fmt.Sprintf("no minimum prices found for symbols: %v", validSymbols))
		return
	}

	writeSuccess(w, prices, map[string]interface{}{
		"count":        len(prices),
		"symbols":      validSymbols,
		"type":         "minimum",
		"requested_at": time.Now().UTC().Format(time.RFC3339),
	})
}

// ========== GET /max ==========

func (h *PriceHandler) GetMaxPrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed",
			"only GET method is supported")
		return
	}

	validSymbols, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	prices, err := h.useCase.GetMaxPrices(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if len(prices) == 0 {
		writeError(w, http.StatusNotFound, "no data",
			fmt.Sprintf("no maximum prices found for symbols: %v", validSymbols))
		return
	}

	writeSuccess(w, prices, map[string]interface{}{
		"count":        len(prices),
		"symbols":      validSymbols,
		"type":         "maximum",
		"requested_at": time.Now().UTC().Format(time.RFC3339),
	})
}

// ========== GET /change ==========

func (h *PriceHandler) GetChangePrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed",
			"only GET method is supported")
		return
	}

	validSymbols, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	changes, err := h.useCase.GetChangePercent(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if len(changes) == 0 {
		writeError(w, http.StatusNotFound, "no data",
			fmt.Sprintf("no change data found for symbols: %v", validSymbols))
		return
	}

	writeSuccess(w, changes, map[string]interface{}{
		"count":        len(changes),
		"symbols":      validSymbols,
		"type":         "change_percent",
		"requested_at": time.Now().UTC().Format(time.RFC3339),
	})
}

// ========== РОУТЕР ==========

type ChiRouter struct {
	mux *chi.Mux
}

func NewChiRouter() *ChiRouter {
	return &ChiRouter{mux: chi.NewRouter()}
}

func (r *ChiRouter) Get(pattern string, handler http.HandlerFunc) {
	r.mux.Get(pattern, handler)
}

func (r *ChiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// ========== ЗАПУСК ==========

func RunServer(uc PriceUseCase) {
	router := NewChiRouter()
	handler, err := NewPriceHandler(uc)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	// Цены
	router.Get("/prices", handler.GetLastPrices)          // GET /prices?symbols=BTC,ETH
	router.Get("/prices/min", handler.GetMinPrices)       // GET /prices/min?symbols=BTC,ETH
	router.Get("/prices/max", handler.GetMaxPrices)       // GET /prices/max?symbols=BTC,ETH
	router.Get("/prices/change", handler.GetChangePrices) // GET /prices/change?symbols=BTC,ETH

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
