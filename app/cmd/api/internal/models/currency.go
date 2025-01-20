package models

type CurrencyConvert struct {
	Amount       float64 `json:"amount" binding:"required"`
	FromCurrency string  `json:"from_currency" binding:"required,len=3"`
	ToCurrency   string  `json:"to_currency" binding:"required,len=3"`
}

type Currency struct {
	Exchange float64 `json:"exchange"`
	Currency string  `json:"currency"`
}

type CurrencyDataResponse struct {
	Rates map[string]float64 `json:"rates"`
}
