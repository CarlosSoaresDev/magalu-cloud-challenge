package models

type Transaction struct {
	ID                string              `json:"gateway"`
	Amount            float64             `json:"amount"`
	Currency          string              `json:"currency"`
	TransactionStatus []TransactionStatus `json:"transaction_status"`
}

type TransactionStatus struct {
	Status   string `json:"status" `
	DateTime string `json:"dateTime"`
}
