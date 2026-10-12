package model

import "time"

type BankAccount struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountNumber string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"account_number"`
	BankName      string    `gorm:"type:varchar(100);not null" json:"bank_name"`
	AccountName   string    `gorm:"type:varchar(100);not null" json:"account_name"`
	DeviceID      string    `gorm:"type:varchar(100)" json:"device_id"`
	DeviceCode    string    `gorm:"type:varchar(100)" json:"device_code"`
	PortAppium    int       `gorm:"default:4723" json:"port_appium"`
	SafeKey       string    `gorm:"type:varchar(100)" json:"safe_key"`
	Pin           string    `gorm:"type:varchar(10)" json:"pin"`
	Username      string    `gorm:"type:varchar(100)" json:"username"`
	Password      string    `gorm:"type:varchar(100)" json:"password"`
	DailyLimit    float64   `gorm:"type:decimal(16,2);default:0" json:"daily_limit"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeviceResponse struct {
	ID         int    `json:"id,omitempty"`
	Code       string `json:"code"`
	PortAppium int    `json:"port_appium"`
}

type BankAccountResponse struct {
	ID            int            `json:"id"`
	AccountNumber string         `json:"account_number"`
	BankName      string         `json:"bank_name"`
	AccountName   string         `json:"account_name"`
	DeviceID      string         `json:"device_id"`
	Device        DeviceResponse `json:"device"`
	SafeKey       string         `json:"safe_key"`
	Pin           string         `json:"pin"`
	Username      string         `json:"username"`
	Password      string         `json:"password"`
}

type DailyLimitResponse struct {
	AccountNumber string  `json:"account_number"`
	DailyLimit    float64 `json:"daily_limit"`
	DailyLimitTrx int     `json:"daily_limit_transaction"`
	CountToday    int     `json:"count_today"`
	IsReached     bool    `json:"is_reached"`
	Reached       bool    `json:"reached"`
}
