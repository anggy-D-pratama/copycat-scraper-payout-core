package usecase

import (
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
	"errors"
	"strconv"
)

type BankTransactionUsecase interface {
	BulkInsert(requests []model.BankTransactionRequest) (model.BankTransactionResponse, error)
}

type bankTransactionUsecase struct {
	repo repository.BankTransactionRepository
}

func NewBankTransactionUsecase(repo repository.BankTransactionRepository) BankTransactionUsecase {
	return &bankTransactionUsecase{repo: repo}
}

func (u *bankTransactionUsecase) BulkInsert(requests []model.BankTransactionRequest) (model.BankTransactionResponse, error) {
	if len(requests) == 0 {
		return model.BankTransactionResponse{}, errors.New("payload array tidak boleh kosong")
	}

	// Validasi atomik: semua item harus valid sebelum insert
	for i, req := range requests {
		if req.BankAccountID <= 0 {
			return model.BankTransactionResponse{}, errors.New("bank_account_id wajib diisi pada item index " + strconv.Itoa(i))
		}
		if req.TransactionDate == "" {
			return model.BankTransactionResponse{}, errors.New("transaction_date wajib diisi pada item index " + strconv.Itoa(i))
		}
		if req.ValueDate == "" {
			return model.BankTransactionResponse{}, errors.New("value_date wajib diisi pada item index " + strconv.Itoa(i))
		}
		if req.Description == "" {
			return model.BankTransactionResponse{}, errors.New("description wajib diisi pada item index " + strconv.Itoa(i))
		}
	}

	// Map DTO ke entity
	transactions := make([]model.BankTransaction, len(requests))
	for i, req := range requests {
		transactions[i] = model.BankTransaction{
			BankAccountID:   req.BankAccountID,
			TransactionDate: req.TransactionDate,
			ValueDate:       req.ValueDate,
			Description:     req.Description,
			Amount:          req.Amount,
			Balance:         req.Balance,
			BeginBalance:    req.BeginBalance,
			IsBalanceOut:    req.IsBalanceOut,
			Source:          req.Source,
			MethodType:      req.MethodType,
			UTR:             req.UTR,
			TransactionCode: req.TransactionCode,
		}
	}

	inserted, err := u.repo.BulkInsert(transactions)
	if err != nil {
		return model.BankTransactionResponse{}, err
	}

	return model.BankTransactionResponse{
		Total:    len(requests),
		Inserted: inserted,
	}, nil
}
