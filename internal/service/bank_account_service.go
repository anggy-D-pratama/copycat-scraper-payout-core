package service

import (
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type BankAccountService interface {
	GetByAccountNumber(accountNumber string) (*model.BankAccountResponse, error)
	GetDailyLimit(accountNumber string) (*model.DailyLimitResponse, error)
}

type bankAccountService struct {
	repo repository.BankAccountRepository
}

func NewBankAccountService(repo repository.BankAccountRepository) BankAccountService {
	return &bankAccountService{repo: repo}
}

func (s *bankAccountService) GetByAccountNumber(accountNumber string) (*model.BankAccountResponse, error) {
	acc, err := s.repo.FindByAccountNumber(accountNumber)
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

func (s *bankAccountService) GetDailyLimit(accountNumber string) (*model.DailyLimitResponse, error) {
	acc, err := s.repo.FindByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}

	return &model.DailyLimitResponse{
		AccountNumber: acc.AccountNumber,
		DailyLimit:    acc.DailyLimit,
	}, nil
}
