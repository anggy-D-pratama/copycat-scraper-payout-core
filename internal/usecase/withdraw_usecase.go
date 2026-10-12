package usecase

import (
	"errors"
	"strings"
	"time"

	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type WithdrawUsecase interface {
	GetWithdraws(bankAccountID int, status string) ([]model.WithdrawResponse, error)
	UpdateWithdraw(req model.WithdrawUpdateRequest) error
	ValidateWithdraw(req model.WithdrawValidateRequest) error
}

type withdrawUsecase struct {
	repo repository.WithdrawRepository
}

func NewWithdrawUsecase(repo repository.WithdrawRepository) WithdrawUsecase {
	return &withdrawUsecase{repo: repo}
}

func (u *withdrawUsecase) GetWithdraws(bankAccountID int, status string) ([]model.WithdrawResponse, error) {
	withdraws, err := u.repo.FindByFilter(bankAccountID, status)
	if err != nil {
		return nil, err
	}

	res := make([]model.WithdrawResponse, 0, len(withdraws))
	for _, w := range withdraws {
		res = append(res, model.WithdrawResponse{
			ID:              w.ID,
			TransactionNo:   w.TransactionNo,
			TransactionCode: w.TransactionCode,
			BankCode:        w.BankCode,
			Amount:          w.Amount,
			Remark:          w.Remark,
			Notes:           w.Notes,
			BankAccount:     w.BankAccount,
			BankAccountId:   w.BankAccountId,
			UniqueCode:      w.UniqueCode,
			Status:          strings.ToLower(w.Status),
			UpdatedAt:       w.UpdatedAt.Format(time.RFC3339),
		})
	}
	return res, nil
}

func (u *withdrawUsecase) UpdateWithdraw(req model.WithdrawUpdateRequest) error {
	w, err := u.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

	w.Status = req.Status
	if req.Notes != "" {
		w.Notes = req.Notes
	}
	if req.UniqueCode != nil {
		w.UniqueCode = req.UniqueCode
	}

	return u.repo.Update(w)
}

func (u *withdrawUsecase) ValidateWithdraw(req model.WithdrawValidateRequest) error {
	w, err := u.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

	if w.TransactionNo != req.TransactionNo {
		return errors.New("invalid transaction_no")
	}

	// Simulasi validasi. Pada environment asli, validate code ini bisa ke service lain.
	// Di mock kita cukup cek transactionCode ada
	if req.TransactionCode == "" {
		return errors.New("transaction_code is required for validation")
	}

	return nil
}
