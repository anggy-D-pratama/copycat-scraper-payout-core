package repository

import (
	"boilerplate/internal/model"
	"gorm.io/gorm"
)

type BankAccountRepository interface {
	FindByAccountNumber(accountNumber string) (*model.BankAccount, error)
}

type bankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) BankAccountRepository {
	return &bankAccountRepository{db: db}
}

func (r *bankAccountRepository) FindByAccountNumber(accountNumber string) (*model.BankAccount, error) {
	var bankAccount model.BankAccount
	err := r.db.Where("account_number = ?", accountNumber).First(&bankAccount).Error
	if err != nil {
		return nil, err
	}
	return &bankAccount, nil
}
