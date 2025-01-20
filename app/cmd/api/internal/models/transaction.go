package models

type Transaction struct {
	Id                string              `json:"id"`
	Amount            float64             `json:"amount"`
	Currency          string              `json:"currency"`
	TransactionStatus []TransactionStatus `json:"transaction_status"`
}

type TransactionStatus struct {
	Status   string `json:"status" `
	DateTime string `json:"dateTime"`
}
