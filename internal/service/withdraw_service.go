package service

import (
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type WithdrawService interface {
	GetWithdraws(bankAccountID int, status string) ([]model.WithdrawResponse, error)
}

type withdrawService struct {
	repo repository.WithdrawRepository
}

func NewWithdrawService(repo repository.WithdrawRepository) WithdrawService {
	return &withdrawService{repo: repo}
}

func (s *withdrawService) GetWithdraws(bankAccountID int, status string) ([]model.WithdrawResponse, error) {
	withdraws, err := s.repo.FindByFilter(bankAccountID, status)
	if err != nil {
		return nil, err
	}

	var res []model.WithdrawResponse
	for _, w := range withdraws {
		res = append(res, model.WithdrawResponse{
			ID:            w.ID,
			TransactionNo: w.TransactionNo,
			Amount:        w.Amount,
			Remark:        w.Remark,
			BankAccount:   w.BankAccount,
			UniqueCode:    w.UniqueCode,
		})
	}
	return res, nil
}
