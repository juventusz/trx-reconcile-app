package model

import "time"

const (
	TransactionTypeDebit  = "DEBIT"
	TransactionTypeCredit = "CREDIT"
)

type Transaction struct {
	TransactionID   string    `json:"trx_id" db:"transaction_id"`
	Amount          float64   `json:"amount" db:"amount"`
	Type            string    `json:"type" db:"type"`
	TransactionTime time.Time `json:"transaction_time" db:"transaction_time"`
}

type BankStatement struct {
	UniqueID string    `json:"unique_identifier"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
}
