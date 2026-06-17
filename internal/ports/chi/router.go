package chi

import (
	"encoding/json"
	"net/http"
	"strings"

	"final/internal/usecases"

	"github.com/go-chi/chi/v5"
)

// ========== РОУТЕР ==========

// Router — интерфейс для HTTP роутера
type Router interface {
	http.Handler
	Get(pattern string, handler http.HandlerFunc)
}

// ChiRouter — реализация Router через chi
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

// ========== ОБРАБОТЧИКИ ==========

type PriceHandler struct {
	useCase usecases.PriceUseCase
}

func NewPriceHandler(uc usecases.PriceUseCase) *PriceHandler {
	return &PriceHandler{useCase: uc}
}

// GET /prices
func (h *PriceHandler) GetPricesLast(w http.ResponseWriter, r *http.Request) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		http.Error(w, "missing symbols param", http.StatusBadRequest)
		return
	}
	symbols := strings.Split(symbolsParam, ",")

	prices, err := h.useCase.GetPricesLast(r.Context(), symbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

// GET /min
func (h *PriceHandler) GetMinPrices(w http.ResponseWriter, r *http.Request) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		http.Error(w, "missing symbols param", http.StatusBadRequest)
		return
	}
	symbols := strings.Split(symbolsParam, ",")

	prices, err := h.useCase.GetMinPrices(r.Context(), symbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

// GET /max
func (h *PriceHandler) GetMaxPrices(w http.ResponseWriter, r *http.Request) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		http.Error(w, "missing symbols param", http.StatusBadRequest)
		return
	}
	symbols := strings.Split(symbolsParam, ",")

	prices, err := h.useCase.GetMaxPrices(r.Context(), symbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

// GET /change
func (h *PriceHandler) GetChangePercent(w http.ResponseWriter, r *http.Request) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		http.Error(w, "missing symbols param", http.StatusBadRequest)
		return
	}
	symbols := strings.Split(symbolsParam, ",")

	changes, err := h.useCase.GetChangePercent(r.Context(), symbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(changes)
}
