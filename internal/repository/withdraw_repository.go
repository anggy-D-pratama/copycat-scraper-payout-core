package repository

import (
	"boilerplate/internal/model"
	"gorm.io/gorm"
)

type WithdrawRepository interface {
	FindByFilter(bankAccountID int, status string) ([]model.Withdraw, error)
}

type withdrawRepository struct {
	db *gorm.DB
}

func NewWithdrawRepository(db *gorm.DB) WithdrawRepository {
	return &withdrawRepository{db: db}
}

func (r *withdrawRepository) FindByFilter(bankAccountID int, status string) ([]model.Withdraw, error) {
	var withdraws []model.Withdraw
	query := r.db.Model(&model.Withdraw{})

	if bankAccountID > 0 {
		query = query.Where("bank_account_id = ?", bankAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&withdraws).Error
	return withdraws, err
}
