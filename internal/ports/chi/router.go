package chi

import (
	"context"
	"encoding/json"
	entity "final/internal/entities"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type ChiRouter struct {
	mux     *chi.Mux
	useCase PriceUseCase
}

func NewChiRouter(uc PriceUseCase) (*ChiRouter, error) {
	apiKey := os.Getenv("COINDESK_API_KEY")
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC&tsyms=USD&api_key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("external API not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	r := &ChiRouter{
		mux:     chi.NewRouter(),
		useCase: uc,
	}
	r.registerRoutes()
	return r, nil
}
func (rt *ChiRouter) registerRoutes() {
	rt.mux.Get("/get/prices/last", rt.GetLastPrices)
	rt.mux.Get("/get/prices/min", rt.GetMinPrices)
	rt.mux.Get("/get/prices/max", rt.GetMaxPrices)
	rt.mux.Get("/get/prices/percent", rt.GetChangePrices)
}

func (r *ChiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
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
func parseAndValidateSymbols(r *http.Request) ([]string, int, error) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("query parameter 'symbols' is required")
	}

	if strings.Trim(symbolsParam, " ,") == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("symbols parameter is empty")
	}

	symbols := strings.Split(symbolsParam, ",")
	validSymbols := make([]string, 0, len(symbols))

	for _, symbol := range symbols {
		symbol = strings.TrimSpace(strings.ToUpper(symbol))
		if symbol == "" {
			continue
		}
		if len(symbol) < 2 || len(symbol) > 10 {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid symbol format: %s", symbol)
		}
		validSymbols = append(validSymbols, symbol)
	}

	if len(validSymbols) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("no valid symbols found")
	}

	return validSymbols, http.StatusOK, nil
}
func toPriceResponse(prices []entity.Price) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(prices))
	for _, p := range prices {
		result = append(result, map[string]interface{}{
			"symbol": p.Symbol,
			"price":  p.Price,
			"time":   p.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
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

// handleServiceError — логирует ошибку и возвращает безопасный ответ клиенту
func handleServiceError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %v", err)                                           // детали для разработчика
	http.Error(w, "internal server error", http.StatusInternalServerError) // безопасно для клиента
}

// ========== GET /prices ==========

func (rt *ChiRouter) GetLastPrices(w http.ResponseWriter, r *http.Request) {
	// 1. Парсинг и валидация symbols
	validSymbols, status, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, status, "invalid request", err.Error())
		return
	}

	// 2. Таймаут контекста из .env
	timeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_SEC"))
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 3. Вызов бизнес-логики
	prices, err := rt.useCase.GetPricesLast(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return

	}
	// 4. Успешный ответ
	response := toPriceResponse(prices)
	json.NewEncoder(w).Encode(response)
}

// ========== GET /min ==========

func (rt *ChiRouter) GetMinPrices(w http.ResponseWriter, r *http.Request) {
	// 1. Парсинг и валидация symbols
	validSymbols, status, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, status, "invalid request", err.Error())
		return
	}

	// 2. Таймаут контекста из .env
	timeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_SEC"))
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 3. Вызов бизнес-логики
	prices, err := rt.useCase.GetMinPrices(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// 4. Успешный ответ
	response := toPriceResponse(prices)
	json.NewEncoder(w).Encode(response)
}

// ========== GET /max ==========

func (rt *ChiRouter) GetMaxPrices(w http.ResponseWriter, r *http.Request) {
	// 1. Парсинг и валидация symbols
	validSymbols, status, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, status, "invalid request", err.Error())
		return
	}
	// 2. Таймаут контекста из .env
	timeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_SEC"))
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 3. Вызов бизнес-логики
	prices, err := rt.useCase.GetMaxPrices(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// 4. Успешный ответ
	response := toPriceResponse(prices)
	json.NewEncoder(w).Encode(response)
}

// ========== GET /change ==========

func (rt *ChiRouter) GetChangePrices(w http.ResponseWriter, r *http.Request) {
	// 1. Парсинг и валидация symbols
	validSymbols, status, err := parseAndValidateSymbols(r)
	if err != nil {
		writeError(w, status, "invalid request", err.Error())
		return
	}

	// 2. Таймаут контекста из .env
	timeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_SEC"))
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 3. Вызов бизнес-логики
	changes, err := rt.useCase.GetChangePercent(ctx, validSymbols)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// 4. Успешный ответ
	response := toPriceResponse(changes)
	json.NewEncoder(w).Encode(response)
}
