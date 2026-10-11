package usecase

import (
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type BankAccountUsecase interface {
	GetByAccountNumber(accountNumber string) (*model.BankAccountResponse, error)
	GetDailyLimit(accountNumber string) (*model.DailyLimitResponse, error)
}

type bankAccountUsecase struct {
	repo repository.BankAccountRepository
}

func NewBankAccountUsecase(repo repository.BankAccountRepository) BankAccountUsecase {
	return &bankAccountUsecase{repo: repo}
}

func (u *bankAccountUsecase) GetByAccountNumber(accountNumber string) (*model.BankAccountResponse, error) {
	acc, err := u.repo.FindByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}

	return &model.BankAccountResponse{
		ID:            acc.ID,
		AccountNumber: acc.AccountNumber,
		BankName:      acc.BankName,
		AccountName:   acc.AccountName,
		DeviceID:      acc.DeviceID,
		Pin:           acc.Pin,
		Username:      acc.Username,
		Password:      acc.Password,
	}, nil
}

func (u *bankAccountUsecase) GetDailyLimit(accountNumber string) (*model.DailyLimitResponse, error) {
	acc, err := u.repo.FindByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}

	return &model.DailyLimitResponse{
		AccountNumber: acc.AccountNumber,
		DailyLimit:    acc.DailyLimit,
	}, nil
}
