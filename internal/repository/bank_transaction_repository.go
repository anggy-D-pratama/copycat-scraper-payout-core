package repository

import (
	"boilerplate/internal/model"

	"gorm.io/gorm"
)

type BankTransactionRepository interface {
	BulkInsert(transactions []model.BankTransaction) (int, error)
}

type bankTransactionRepository struct {
	db *gorm.DB
}

func NewBankTransactionRepository(db *gorm.DB) BankTransactionRepository {
	return &bankTransactionRepository{db: db}
}

func (r *bankTransactionRepository) BulkInsert(transactions []model.BankTransaction) (int, error) {
	if len(transactions) == 0 {
		return 0, nil
	}

	inserted := 0

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, trx := range transactions {
			var exists int64
			if trx.UTR != "" || trx.TransactionCode != "" {
				// Outgoing deduplication: by bank_account_id AND (utr OR transaction_code)
				if trx.UTR != "" && trx.TransactionCode != "" {
					if err := tx.Model(&model.BankTransaction{}).
						Where("bank_account_id = ? AND (utr = ? OR transaction_code = ?)", trx.BankAccountID, trx.UTR, trx.TransactionCode).
						Count(&exists).Error; err != nil {
						return err
					}
				} else if trx.UTR != "" {
					if err := tx.Model(&model.BankTransaction{}).
						Where("bank_account_id = ? AND utr = ?", trx.BankAccountID, trx.UTR).
						Count(&exists).Error; err != nil {
						return err
					}
				} else {
					if err := tx.Model(&model.BankTransaction{}).
						Where("bank_account_id = ? AND transaction_code = ?", trx.BankAccountID, trx.TransactionCode).
						Count(&exists).Error; err != nil {
						return err
					}
				}
			} else {
				// Incoming deduplication: by bank_account_id AND description
				if err := tx.Model(&model.BankTransaction{}).
					Where("bank_account_id = ? AND description = ?", trx.BankAccountID, trx.Description).
					Count(&exists).Error; err != nil {
					return err
				}
			}

			if exists > 0 {
				continue
			}

			if err := tx.Create(&trx).Error; err != nil {
				return err
			}
			inserted++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return inserted, nil
}
