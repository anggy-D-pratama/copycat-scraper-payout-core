package model

import "time"

type BankTransaction struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	BankAccountID   int       `gorm:"index;not null" json:"bank_account_id"`
	TransactionDate string    `gorm:"type:varchar(10);not null" json:"transaction_date"`
	ValueDate       string    `gorm:"type:varchar(10);not null" json:"value_date"`
	Description     string    `gorm:"type:text;not null" json:"description"`
	Amount          float64   `gorm:"type:decimal(16,2);not null" json:"amount"`
	Balance         float64   `gorm:"type:decimal(16,2);not null" json:"balance"`
	BeginBalance    float64   `gorm:"type:decimal(16,2);not null" json:"begin_balance"`
	IsBalanceOut    bool      `gorm:"not null" json:"is_balance_out"`
	Source          string    `gorm:"type:varchar(50)" json:"source"`
	MethodType      string    `gorm:"type:varchar(50)" json:"method_type"`
	UTR             string    `gorm:"type:varchar(100)" json:"utr"`
	TransactionCode string    `gorm:"type:varchar(100)" json:"transaction_code"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BankTransactionRequest struct {
	TransactionDate string  `json:"transaction_date"`
	ValueDate       string  `json:"value_date"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Balance         float64 `json:"balance"`
	BeginBalance    float64 `json:"begin_balance"`
	IsBalanceOut    bool    `json:"is_balance_out"`
	Source          string  `json:"source"`
	MethodType      string  `json:"method_type"`
	UTR             string  `json:"utr"`
	TransactionCode string  `json:"transaction_code"`
	BankAccountID   int     `json:"bank_account_id"`
}

type BankTransactionResponse struct {
	Total    int `json:"total"`
	Inserted int `json:"inserted"`
}
