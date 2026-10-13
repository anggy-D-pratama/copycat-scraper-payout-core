package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"boilerplate/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockBankTransactionUsecase struct{}

func (m *mockBankTransactionUsecase) BulkInsert(requests []model.BankTransactionRequest) (model.BankTransactionResponse, error) {
	// Simulate validation failure for test cases
	for i, req := range requests {
		if req.BankAccountID <= 0 {
			return model.BankTransactionResponse{}, errors.New("bank_account_id wajib diisi pada item index " + fmt.Sprint(i))
		}
		if req.TransactionDate == "" {
			return model.BankTransactionResponse{}, errors.New("transaction_date wajib diisi pada item index " + fmt.Sprint(i))
		}
		if req.ValueDate == "" {
			return model.BankTransactionResponse{}, errors.New("value_date wajib diisi pada item index " + fmt.Sprint(i))
		}
		if req.Description == "" {
			return model.BankTransactionResponse{}, errors.New("description wajib diisi pada item index " + fmt.Sprint(i))
		}
	}

	// Return mock response
	return model.BankTransactionResponse{
		Total:    len(requests),
		Inserted: len(requests),
	}, nil
}

func TestBankTransactionHandler_BulkInsert(t *testing.T) {
	apiKey := "test-key"
	secretKey := "test-secret"
	os.Setenv("API_KEY", apiKey)
	os.Setenv("SECRET_KEY", secretKey)
	defer func() {
		os.Unsetenv("API_KEY")
		os.Unsetenv("SECRET_KEY")
	}()

	mockUsecase := &mockBankTransactionUsecase{}
	handler := NewBankTransactionHandler(mockUsecase)

	t.Run("success bulk insert", func(t *testing.T) {
		app := fiber.New()
		app.Post("/bank-transactions/bulk", handler.BulkInsert)

		payload := []model.BankTransactionRequest{
			{
				BankAccountID:   1,
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "Transfer masuk dari JOHN DOE",
				Amount:          100000,
				Balance:         500000,
				BeginBalance:    400000,
				IsBalanceOut:    false,
				Source:          "bank-scraper",
				MethodType:      "rtgs",
				UTR:             "UTR001",
			},
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/bank-transactions/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		timestamp := fmt.Sprint(time.Now().Unix())
		sig := genSig(secretKey, apiKey, timestamp)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "success", result["status"])
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["total"])
		assert.Equal(t, float64(1), data["inserted"])
	})

	t.Run("empty payload array", func(t *testing.T) {
		app := fiber.New()
		app.Post("/bank-transactions/bulk", handler.BulkInsert)

		payload := []model.BankTransactionRequest{}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/bank-transactions/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		timestamp := fmt.Sprint(time.Now().Unix())
		sig := genSig(secretKey, apiKey, timestamp)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "error", result["status"])
		assert.Contains(t, result["data"], "tidak boleh kosong")
	})

	t.Run("invalid request body", func(t *testing.T) {
		app := fiber.New()
		app.Post("/bank-transactions/bulk", handler.BulkInsert)

		req := httptest.NewRequest("POST", "/bank-transactions/bulk", bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		timestamp := fmt.Sprint(time.Now().Unix())
		sig := genSig(secretKey, apiKey, timestamp)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "error", result["status"])
		assert.Contains(t, result["data"], "invalid request body")
	})

	t.Run("validation error - missing bank_account_id", func(t *testing.T) {
		app := fiber.New()
		app.Post("/bank-transactions/bulk", handler.BulkInsert)

		payload := []model.BankTransactionRequest{
			{
				BankAccountID:   0, // invalid
				TransactionDate: "2026-10-13",
				ValueDate:       "2026-10-13",
				Description:     "Transfer masuk",
				Amount:          100000,
			},
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/bank-transactions/bulk", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		timestamp := fmt.Sprint(time.Now().Unix())
		sig := genSig(secretKey, apiKey, timestamp)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "error", result["status"])
		assert.Contains(t, result["data"], "bank_account_id wajib diisi")
	})
}
