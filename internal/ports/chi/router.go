package chi

import (
	"encoding/json"
	dto "final/pkg"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router  *chi.Mux
	service PriceUseCase
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
func NewServer(service PriceUseCase) (*Server, error) {
	if service == nil { //
		return nil, fmt.Errorf("errprs")
	}
	r := chi.NewRouter()

	return &Server{router: r,
		service: service,
	}, nil
}

func (s *Server) Start() {
	s.router.Get("/api/v1/get/prices/last", s.GetLastPrice)
	s.router.Get("/api/v1/get/prices/min", s.GetMinPrice)
	s.router.Get("/api/v1/get/prices/max", s.GetMaxPrice)
	s.router.Get("/api/v1/get/prices/percent", s.GetChangePercent)

	// Swagger UI
	s.router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Swagger JSON
	s.router.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}

// ================GetLastPrice================//
// @Summary Последние цены
// @Param symbols query string true "Символы через запятую"
// @Success 200 {array} dto.PriceDTO
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/get/prices/last [get]
func (s *Server) GetLastPrice(rw http.ResponseWriter, req *http.Request) {
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}

	splitSymbols := strings.Split(symbols, ",")

	var validSymbols []string
	for _, s := range splitSymbols {
		s = strings.TrimSpace(strings.ToUpper(s))
		if s == "" {
			continue
		}
		if len(s) < 2 || len(s) > 10 {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		valid := true
		for _, c := range s {
			if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				valid = false
				break
			}
		}
		if !valid {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		validSymbols = append(validSymbols, s)
	}
	if len(validSymbols) == 0 {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	if len(validSymbols) > 20 {
		http.Error(rw, "too many symbols, max 20", http.StatusBadRequest)
		return
	}

	prices, err := s.service.GetPricesLast(req.Context(), validSymbols)
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
			Time:   v.CreatedAt.Add(3 * time.Hour).Format("2006-01-02T15:04:05"),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		fmt.Errorf("ERROR: failed to encode response: %v", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// ================GetMaxPrice================//
// @Summary Максимальные цены
// @Param symbols query string true "Символы через запятую"
// @Success 200 {array} dto.PriceDTO
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/get/prices/max [get]
func (s *Server) GetMaxPrice(rw http.ResponseWriter, req *http.Request) {
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}

	splitSymbols := strings.Split(symbols, ",")

	var validSymbols []string
	for _, s := range splitSymbols {
		s = strings.TrimSpace(strings.ToUpper(s))
		if s == "" {
			continue
		}
		if len(s) < 2 || len(s) > 10 {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		valid := true
		for _, c := range s {
			if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				valid = false
				break
			}
		}
		if !valid {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		validSymbols = append(validSymbols, s)
	}
	if len(validSymbols) == 0 {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	if len(validSymbols) > 20 {
		http.Error(rw, "too many symbols, max 20", http.StatusBadRequest)
		return
	}

	prices, err := s.service.GetMaxPrices(req.Context(), validSymbols)
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
			Time:   v.CreatedAt.Add(3 * time.Hour).Format("2006-01-02T15:04:05"),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		fmt.Errorf("ERROR: failed to encode response: %v", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// ================GetMinPrice================//
// @Summary Минимальные цены
// @Param symbols query string true "Символы через запятую"
// @Success 200 {array} dto.PriceDTO
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/get/prices/min [get]
func (s *Server) GetMinPrice(rw http.ResponseWriter, req *http.Request) {
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}

	splitSymbols := strings.Split(symbols, ",")

	var validSymbols []string
	for _, s := range splitSymbols {
		s = strings.TrimSpace(strings.ToUpper(s))
		if s == "" {
			continue
		}
		if len(s) < 2 || len(s) > 10 {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		valid := true
		for _, c := range s {
			if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				valid = false
				break
			}
		}
		if !valid {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		validSymbols = append(validSymbols, s)
	}
	if len(validSymbols) == 0 {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	if len(validSymbols) > 20 {
		http.Error(rw, "too many symbols, max 20", http.StatusBadRequest)
		return
	}

	prices, err := s.service.GetMinPrices(req.Context(), validSymbols)
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
			Time:   v.CreatedAt.Add(3 * time.Hour).Format("2006-01-02T15:04:05"),
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		fmt.Errorf("ERROR: failed to encode response: %v", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// ================GetChangePrices================//
// @Summary Изменение в процентах
// @Param symbols query string true "Символы через запятую"
// @Success 200 {array} dto.PriceDTO
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/get/prices/percent [get]
func (s *Server) GetChangePercent(rw http.ResponseWriter, req *http.Request) {
	symbols := req.URL.Query().Get("symbols")
	if symbols == "" {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}

	splitSymbols := strings.Split(symbols, ",")

	var validSymbols []string
	for _, s := range splitSymbols {
		s = strings.TrimSpace(strings.ToUpper(s))
		if s == "" {
			continue
		}
		if len(s) < 2 || len(s) > 10 {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		valid := true
		for _, c := range s {
			if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				valid = false
				break
			}
		}
		if !valid {
			http.Error(rw, fmt.Sprintf("invalid symbol: %s", s), http.StatusBadRequest)
			return
		}
		validSymbols = append(validSymbols, s)
	}
	if len(validSymbols) == 0 {
		http.Error(rw, "symbols is required", http.StatusBadRequest)
		return
	}
	if len(validSymbols) > 20 {
		http.Error(rw, "too many symbols, max 20", http.StatusBadRequest)
		return
	}

	prices, err := s.service.GetChangePercent(req.Context(), validSymbols)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "no data for this currency yet, try again in 5 minutes", http.StatusOK)
		return
	}

	if len(prices) == 0 {
		http.Error(rw, "no data for this currency yet, try again in 5 minutes", http.StatusOK)
		return
	}

	var data []dto.PriceDTO
	found := make(map[string]bool)

	for _, v := range prices {
		found[v.Symbol] = true
		data = append(data, dto.PriceDTO{
			Symbol: v.Symbol,
			Price:  v.Price,
			Time:   v.CreatedAt.Add(3 * time.Hour).Format("2006-01-02T15:04:05"),
		})
	}

	// Проверяем отсутствующие символы
	var missing []string
	for _, s := range validSymbols {
		if !found[s] {
			missing = append(missing, s)
		}
	}

	// Добавляем отсутствующие валюты с сообщением
	for _, s := range missing {
		data = append(data, dto.PriceDTO{
			Symbol:  s,
			Message: "no data yet, try again in 5 minutes",
		})
	}

	rw.Header().Add("ContentType", "application/json")
	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		fmt.Errorf("ERROR: failed to encode response: %v", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
