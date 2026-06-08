package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	entity "final/internal/entities"
	"final/internal/ports/input"
)

type PriceHandler struct {
	service input.PriceService
}

func NewPriceHandler(service input.PriceService) *PriceHandler {
	return &PriceHandler{service: service}
}

func (h *PriceHandler) RegisterRoutes(router Router) {
	router.Get("/prices", h.getPricesLast)
	router.Get("/min", h.getMinPrices)
	router.Get("/max", h.getMaxPrices)
	router.Get("/change", h.getChangePercent)
}

func (h *PriceHandler) getPricesLast(w http.ResponseWriter, r *http.Request) {
	h.respondWithPrices(w, r, h.service.GetPricesLast)
}

func (h *PriceHandler) getMinPrices(w http.ResponseWriter, r *http.Request) {
	h.respondWithPrices(w, r, h.service.GetMinPrices)
}

func (h *PriceHandler) getMaxPrices(w http.ResponseWriter, r *http.Request) {
	h.respondWithPrices(w, r, h.service.GetMaxPrices)
}

func (h *PriceHandler) getChangePercent(w http.ResponseWriter, r *http.Request) {
	h.respondWithPrices(w, r, h.service.GetChangePercent)
}

type priceFetcher func(ctx context.Context, symbols []string) ([]entity.Price, error)

func (h *PriceHandler) respondWithPrices(w http.ResponseWriter, r *http.Request, fetch priceFetcher) {
	symbols, ok := parseSymbols(w, r)
	if !ok {
		return
	}

	prices, err := fetch(r.Context(), symbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

func parseSymbols(w http.ResponseWriter, r *http.Request) ([]string, bool) {
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		http.Error(w, "missing symbols param", http.StatusBadRequest)
		return nil, false
	}
	return strings.Split(symbolsParam, ","), true
}
