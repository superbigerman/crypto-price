package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"final/internal/config"
	"final/internal/external"
)

func main() {
	cfg := config.Load()
	client := external.NewCoinDeskClient(cfg)

	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		symbolsParam := r.URL.Query().Get("symbols")
		if symbolsParam == "" {
			http.Error(w, "missing symbols param", http.StatusBadRequest)
			return
		}
		symbols := strings.Split(symbolsParam, ",")

		prices, err := client.GetRealTimePrices(r.Context(), symbols)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
