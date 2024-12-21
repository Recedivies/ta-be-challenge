package models

type Transaction struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}
