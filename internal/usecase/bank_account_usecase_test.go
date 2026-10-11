package usecase_test

import (
	"errors"
	"testing"

	"boilerplate/internal/model"
	"boilerplate/internal/usecase"
)

type mockBankAccountRepo struct {
	data map[string]*model.BankAccount
}

func (m *mockBankAccountRepo) FindByAccountNumber(accNo string) (*model.BankAccount, error) {
	if val, ok := m.data[accNo]; ok {
		return val, nil
	}
	return nil, errors.New("not found")
}

func TestBankAccountUsecase(t *testing.T) {
	repo := &mockBankAccountRepo{
		data: map[string]*model.BankAccount{
			"12345": {
				ID:            1,
				AccountNumber: "12345",
				BankName:      "BCA",
				AccountName:   "Tester",
				DeviceID:      "dev-1",
				DeviceCode:    "DEV-01",
				PortAppium:    4723,
				SafeKey:       "654321",
				Pin:           "654321",
				Username:      "tester",
				Password:      "secret",
				DailyLimit:    5000000,
			},
		},
	}
	svc := usecase.NewBankAccountUsecase(repo)

	// Test GetByAccountNumber Success
	acc, err := svc.GetByAccountNumber("12345")
	if err != nil || acc.AccountNumber != "12345" {
		t.Errorf("expected acc, got error or mismatch: %v", err)
	}
	if acc.Device.PortAppium != 4723 || acc.Device.Code != "DEV-01" || acc.SafeKey != "654321" {
		t.Errorf("expected nested device and safe_key, got %+v", acc)
	}

	// Test GetByAccountNumber Not Found
	_, err = svc.GetByAccountNumber("99999")
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Test GetDailyLimit Success
	limit, err := svc.GetDailyLimit("12345")
	if err != nil || limit.DailyLimit != 5000000 {
		t.Errorf("expected limit 5000000, got: %v", limit)
	}
	if limit.IsReached != false || limit.Reached != false {
		t.Errorf("expected is_reached and reached to be false, got: %+v", limit)
	}
}
