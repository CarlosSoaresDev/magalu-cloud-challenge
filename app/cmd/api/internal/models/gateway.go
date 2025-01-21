package models

type CardDetails struct {
	Number string `json:"number" binding:"required,cnumber"`
	Expiry string `json:"expiry" binding:"required,cexpirate"`
	Cvv    string `json:"cvv" binding:"required,len=3"`
}

type Gateway struct {
	Gateway       string      `json:"gateway" binding:"required"`
	Amount        float64     `json:"amount" binding:"required"`
	Currency      string      `json:"currency" binding:"required,len=3"`
	PaymentMethod string      `json:"payment_method" binding:"required"`
	CardDetails   CardDetails `json:"card_details" binding:"required"`
}
