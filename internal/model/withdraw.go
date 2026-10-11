package model

import "time"

type Withdraw struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionNo   string    `gorm:"type:varchar(100);not null" json:"transaction_no"`
	TransactionCode string    `gorm:"type:varchar(100)" json:"transaction_code"`
	BankCode        string    `gorm:"type:varchar(50)" json:"bank_code"`
	Amount          float64   `gorm:"type:decimal(16,2);not null" json:"amount"`
	Remark          string    `gorm:"type:text" json:"remark"`
	Notes           string    `gorm:"type:text" json:"notes"`
	BankAccount     string    `gorm:"type:varchar(100);not null" json:"bank_account"`
	BankAccountId   int       `gorm:"not null;index" json:"bank_account_id"`
	UniqueCode      *string   `gorm:"type:varchar(50)" json:"unique_code"`
	Status          string    `gorm:"type:varchar(50);not null;default:'QUEUE'" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WithdrawResponse struct {
	ID              int     `json:"id"`
	TransactionNo   string  `json:"transaction_no"`
	TransactionCode string  `json:"transaction_code"`
	BankCode        string  `json:"bank_code"`
	Amount          float64 `json:"amount"`
	Remark          string  `json:"remark"`
	BankAccount     string  `json:"bank_account"`
	UniqueCode      *string `json:"unique_code"`
	Status          string  `json:"status"`
}

type WithdrawUpdateRequest struct {
	ID         int     `json:"id" validate:"required"`
	Status     string  `json:"status" validate:"required"`
	Notes      string  `json:"notes"`
	UniqueCode *string `json:"unique_code"`
}

type WithdrawValidateRequest struct {
	ID              int    `json:"id" validate:"required"`
	TransactionNo   string `json:"transaction_no" validate:"required"`
	TransactionCode string `json:"transaction_code" validate:"required"`
}
