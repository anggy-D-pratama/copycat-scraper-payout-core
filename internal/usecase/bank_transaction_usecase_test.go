package usecase

import (
	"boilerplate/internal/model"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBankTransactionRepository struct {
	shouldReturnError bool
}

func (m *mockBankTransactionRepository) BulkInsert(transactions []model.BankTransaction) (int, error) {
	if m.shouldReturnError {
		return 0, errors.New("database error")
	}
	return len(transactions), nil
}

func TestBankTransactionUsecase_BulkInsert(t *testing.T) {
	t.Run("success bulk insert", func(t *testing.T) {
		repo := &mockBankTransactionRepository{shouldReturnError: false}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{
			{
				BankAccountID:   1,
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "Transfer masuk dari John Doe",
				Amount:          100000,
				Balance:         500000,
				BeginBalance:    400000,
				IsBalanceOut:    false,
				UTR:             "UTR001",
			},
			{
				BankAccountID:   1,
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "Transfer keluar ke Jane Doe",
				Amount:          50000,
				Balance:         450000,
				BeginBalance:    500000,
				IsBalanceOut:    true,
				TransactionCode: "TXC002",
			},
		}

		res, err := usecase.BulkInsert(requests)

		assert.NoError(t, err)
		assert.Equal(t, 2, res.Total)
		assert.Equal(t, 2, res.Inserted)
	})

	t.Run("empty array validation", func(t *testing.T) {
		repo := &mockBankTransactionRepository{}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{}

		_, err := usecase.BulkInsert(requests)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak boleh kosong")
	})

	t.Run("validation error missing bank_account_id", func(t *testing.T) {
		repo := &mockBankTransactionRepository{}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{
			{
				BankAccountID:   0,
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "Test",
			},
		}

		_, err := usecase.BulkInsert(requests)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bank_account_id wajib diisi")
	})

	t.Run("validation error missing transaction_date", func(t *testing.T) {
		repo := &mockBankTransactionRepository{}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{
			{
				BankAccountID:   1,
				TransactionDate: "",
				ValueDate:       "2026-10-13",
				Description:     "Test",
			},
		}

		_, err := usecase.BulkInsert(requests)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction_date wajib diisi")
	})

	t.Run("validation error missing value_date", func(t *testing.T) {
		repo := &mockBankTransactionRepository{}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{
			{
				BankAccountID:   1,
				TransactionDate: "2026-10-13",
				ValueDate:       "",
				Description:     "Test",
			},
		}

		_, err := usecase.BulkInsert(requests)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value_date wajib diisi")
	})

	t.Run("validation error missing description", func(t *testing.T) {
		repo := &mockBankTransactionRepository{}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{
			{
				BankAccountID:   1,
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "",
			},
		}

		_, err := usecase.BulkInsert(requests)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description wajib diisi")
	})

	t.Run("repo error handling", func(t *testing.T) {
		repo := &mockBankTransactionRepository{shouldReturnError: true}
		usecase := NewBankTransactionUsecase(repo)

		requests := []model.BankTransactionRequest{
			{
				BankAccountID:   1,
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "Test",
			},
		}

		_, err := usecase.BulkInsert(requests)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}
