package chi

import (
	"encoding/json"
	dto "final/internal/usecases"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router  *chi.Mux
	service PriceUseCase
}

func NewServer(service PriceUseCase) (*Server, error) {
	if service == nil { //
		return nil, fmt.Errorf("errprs") // хиитрую ошибку добавить
	}
	r := chi.NewRouter()
	s := &Server{router: r, service: service}

	return &Server{router: r,
		service: service,
	}, nil
}

func (s *Server) Start() {
	s.router.Get("/get/prices/last", s.GetLastPrice)
	s.router.Get("/get/prices/min", s.GetMinPrice)
	s.router.Get("/get/prices/max", s.GetMaxPrice)
	s.router.Get("/get/prices/percent", s.GetChangePercent)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}

// ================GetLastPrice================//
func (s *Server) GetLastPrice(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	splitSymbols := strings.Split(symbols, ",")
	prices, err := s.service.GetPricesLast(req.Context(), splitSymbols)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "internal error", http.StatusInternalServerError)
		return
	}
	if len(prices) == 0 {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	var data []dto.PriceDTO
	for _, v := range prices {
		data = append(data, dto.PriceDTO{
			Symbol: v.Symbol,
			Price:  v.Price,
			Time:   v.CreatedAt.Format(time.RFC3339),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	rw.WriteHeader(http.StatusOK)
}

// ================GetMaxPrice================//
func (s *Server) GetMaxPrice(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	splitSymbols := strings.Split(symbols, ",")
	prices, err := s.service.GetMaxPrices(req.Context(), splitSymbols)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "internal error", http.StatusInternalServerError)
		return
	}
	if len(prices) == 0 {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	var data []dto.PriceDTO
	for _, v := range prices {
		data = append(data, dto.PriceDTO{
			Symbol: v.Symbol,
			Price:  v.Price,
			Time:   v.CreatedAt.Format(time.RFC3339),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	rw.WriteHeader(http.StatusOK)
}

// ================GetMinPrice================//
func (s *Server) GetMinPrice(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	splitSymbols := strings.Split(symbols, ",")
	prices, err := s.service.GetMinPrices(req.Context(), splitSymbols)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "internal error", http.StatusInternalServerError)
		return
	}
	if len(prices) == 0 {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	var data []dto.PriceDTO
	for _, v := range prices {
		data = append(data, dto.PriceDTO{
			Symbol: v.Symbol,
			Price:  v.Price,
			Time:   v.CreatedAt.Format(time.RFC3339),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	rw.WriteHeader(http.StatusOK)
}

// ================GetChangePrices================//
func (s *Server) GetChangePercent(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	splitSymbols := strings.Split(symbols, ",")
	prices, err := s.service.GetChangePercent(req.Context(), splitSymbols)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "internal error", http.StatusInternalServerError)
		return
	}
	if len(prices) == 0 {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	var data []dto.PriceDTO
	for _, v := range prices {
		data = append(data, dto.PriceDTO{
			Symbol: v.Symbol,
			Price:  v.Price,
			Time:   v.CreatedAt.Format(time.RFC3339),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	rw.WriteHeader(http.StatusOK)
}
