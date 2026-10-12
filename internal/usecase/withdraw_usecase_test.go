package usecase_test

import (
	"errors"
	"testing"
	"time"

	"boilerplate/internal/model"
	"boilerplate/internal/usecase"
)

type mockWithdrawRepo struct {
	data map[int]*model.Withdraw
}

func (m *mockWithdrawRepo) FindByFilter(accID int, status string) ([]model.Withdraw, error) {
	var res []model.Withdraw
	for _, v := range m.data {
		if (accID == 0 || v.BankAccountId == accID) && (status == "" || v.Status == status) {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (m *mockWithdrawRepo) FindByID(id int) (*model.Withdraw, error) {
	if val, ok := m.data[id]; ok {
		return val, nil
	}
	return nil, errors.New("not found")
}

func (m *mockWithdrawRepo) Update(w *model.Withdraw) error {
	m.data[w.ID] = w
	return nil
}

func TestWithdrawUsecase(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	repo := &mockWithdrawRepo{
		data: map[int]*model.Withdraw{
			1: {ID: 1, BankAccountId: 10, Status: "QUEUE", TransactionNo: "TRX-1", TransactionCode: "TXC-1", BankCode: "BCA", UpdatedAt: now},
		},
	}
	svc := usecase.NewWithdrawUsecase(repo)

	// Test GetWithdraws
	ws, err := svc.GetWithdraws(10, "QUEUE")
	if err != nil || len(ws) != 1 {
		t.Errorf("expected 1 withdraw, got: %v", len(ws))
	}
	if ws[0].BankCode != "BCA" || ws[0].TransactionCode != "TXC-1" {
		t.Errorf("expected BankCode BCA and TransactionCode TXC-1, got %s and %s", ws[0].BankCode, ws[0].TransactionCode)
	}
	if ws[0].UpdatedAt != now.Format(time.RFC3339) {
		t.Errorf("expected UpdatedAt %s, got %s", now.Format(time.RFC3339), ws[0].UpdatedAt)
	}

	// Test UpdateWithdraw
	code := "UCODE123"
	err = svc.UpdateWithdraw(model.WithdrawUpdateRequest{
		ID:         1,
		Status:     "PROCESS",
		UniqueCode: &code,
	})
	if err != nil {
		t.Errorf("update error: %v", err)
	}
	if repo.data[1].Status != "PROCESS" || *repo.data[1].UniqueCode != code {
		t.Errorf("expected update to save in repo")
	}

	// Test ValidateWithdraw - Valid
	err = svc.ValidateWithdraw(model.WithdrawValidateRequest{
		ID:              1,
		TransactionNo:   "TRX-1",
		TransactionCode: "VALIDCODE",
	})
	if err != nil {
		t.Errorf("expected valid, got: %v", err)
	}

	// Test ValidateWithdraw - Invalid TRX No
	err = svc.ValidateWithdraw(model.WithdrawValidateRequest{
		ID:              1,
		TransactionNo:   "TRX-2",
		TransactionCode: "VALIDCODE",
	})
	if err == nil {
		t.Errorf("expected invalid, got nil")
	}
}
