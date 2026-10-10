package model

import "time"

type Withdraw struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionNo string    `gorm:"type:varchar(100);not null" json:"transaction_no"`
	Amount        float64   `gorm:"type:decimal(16,2);not null" json:"amount"`
	Remark        string    `gorm:"type:text" json:"remark"`
	BankAccount   string    `gorm:"type:varchar(100);not null" json:"bank_account"`
	BankAccountId int       `gorm:"not null;index" json:"bank_account_id"`
	UniqueCode    *string   `gorm:"type:varchar(50)" json:"unique_code"`
	Status        string    `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WithdrawResponse struct {
	ID            int     `json:"id"`
	TransactionNo string  `json:"transaction_no"`
	Amount        float64 `json:"amount"`
	Remark        string  `json:"remark"`
	BankAccount   string  `json:"bank_account"`
	UniqueCode    *string `json:"unique_code"`
}
