package usecases

type PriceDTO struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}
