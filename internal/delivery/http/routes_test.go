package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"boilerplate/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockBankAccountUsecase struct{}

func (m *mockBankAccountUsecase) GetByAccountNumber(accountNumber string) (*model.BankAccountResponse, error) {
	if accountNumber == "1234567890" {
		return &model.BankAccountResponse{
			ID:            1,
			AccountNumber: "1234567890",
			BankName:      "BCA",
			AccountName:   "Bot Testing Account",
			DeviceID:      "device-123",
			Device: model.DeviceResponse{
				ID:         1,
				Code:       "device-123",
				PortAppium: 4723,
			},
			SafeKey:  "123456",
			Pin:      "123456",
			Username: "bot_user",
			Password: "password",
		}, nil
	}
	return nil, errors.New("not found")
}

func (m *mockBankAccountUsecase) GetDailyLimit(accountNumber string) (*model.DailyLimitResponse, error) {
	if accountNumber == "1234567890" {
		return &model.DailyLimitResponse{
			AccountNumber: "1234567890",
			DailyLimit:    100000000,
			DailyLimitTrx: 500,
			CountToday:    0,
			IsReached:     false,
			Reached:       false,
		}, nil
	}
	return nil, errors.New("not found")
}

type mockWithdrawUsecase struct{}

func (m *mockWithdrawUsecase) GetWithdraws(bankAccountID int, status string) ([]model.WithdrawResponse, error) {
	return []model.WithdrawResponse{
		{
			ID:              1,
			TransactionNo:   "TRX-001",
			TransactionCode: "TXC-001",
			BankCode:        "BCA",
			Amount:          50000,
			Remark:          "Test withdraw 1",
			BankAccount:     "1234567890",
			BankAccountId:   1,
			Status:          "queue",
			UpdatedAt:       "2026-10-12T10:30:00Z",
		},
	}, nil
}

func (m *mockWithdrawUsecase) UpdateWithdraw(req model.WithdrawUpdateRequest) error {
	if req.ID == 999 {
		return errors.New("not found")
	}
	return nil
}

func (m *mockWithdrawUsecase) ValidateWithdraw(req model.WithdrawValidateRequest) error {
	if req.TransactionNo != "TRX-001" {
		return errors.New("invalid transaction_no")
	}
	return nil
}

func genSig(secretKey, apiKey, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(fmt.Sprintf("%s:%s", apiKey, timestamp)))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestSetupRoutes_HMACAuth(t *testing.T) {
	apiKey := "test-key"
	secretKey := "test-secret"
	os.Setenv("API_KEY", apiKey)
	os.Setenv("SECRET_KEY", secretKey)
	defer func() {
		os.Unsetenv("API_KEY")
		os.Unsetenv("SECRET_KEY")
	}()

	app := fiber.New()

	handlers := &Handlers{
		User:        &UserHandler{},
		Withdraw:    NewWithdrawHandler(&mockWithdrawUsecase{}),
		BankAccount: NewBankAccountHandler(&mockBankAccountUsecase{}),
		Telegram:    NewTelegramHandler(),
	}

	SetupRoutes(app, handlers)

	t.Run("Public Route Health Check - No Auth Header Required", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Swagger Route /swagger/doc.json - Returns 200 and Valid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
	})

	t.Run("Swagger Route /swagger/index.html - Returns 200 HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	})

	t.Run("Protected Route /bank-account - Rejected without Auth Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bank-account", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Protected Route /bank-account/daily-limit - Rejected without Auth Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bank-account/daily-limit", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Protected Route /withdraw - Rejected with Invalid Signature", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/withdraw", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", "123456789")
		req.Header.Set("x-signature", "invalid-sig")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Protected Route /bank-account - Accepted with Valid HMAC Signature", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		req := httptest.NewRequest(http.MethodGet, "/bank-account?account_number=1234567890", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
		data, ok := body["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "1234567890", data["account_number"])
		assert.Equal(t, "123456", data["safe_key"])
		device, ok := data["device"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(4723), device["port_appium"])
	})

	t.Run("Protected Route /bank-account/daily-limit - Valid Envelope Response", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		req := httptest.NewRequest(http.MethodGet, "/bank-account/daily-limit?account_number=1234567890", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
		data, ok := body["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, false, data["is_reached"])
		assert.Equal(t, false, data["reached"])
	})

	t.Run("Protected Route GET /withdraw - Valid Envelope Response", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		req := httptest.NewRequest(http.MethodGet, "/withdraw?bank_account_id=1&status=queue", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
		data, ok := body["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1)
		item, ok := data[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "2026-10-12T10:30:00Z", item["updated_at"])
	})

	t.Run("Protected Route POST /withdraw/validate - Valid Envelope Response", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		payload := []byte(`{"id":1,"transaction_no":"TRX-001","transaction_code":"TXC-001"}`)
		req := httptest.NewRequest(http.MethodPost, "/withdraw/validate", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
	})

	t.Run("Protected Route POST /withdraw/update - Valid Envelope Response", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		payload := []byte(`{"id":1,"status":"process","notes":"ok","unique_code":"REF1"}`)
		req := httptest.NewRequest(http.MethodPost, "/withdraw/update", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
	})

	t.Run("Protected Route POST /telegram/send - Valid Success Response", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		payload := []byte(`{"account_number":"1234567890","text":"Testing bot message","is_photo":"false","is_file":"false"}`)
		req := httptest.NewRequest(http.MethodPost, "/telegram/send", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
	})

	t.Run("Protected Route POST /telegram/save-image - Valid Success Response", func(t *testing.T) {
		ts := "1773223200000"
		sig := genSig(secretKey, apiKey, ts)

		req := httptest.NewRequest(http.MethodPost, "/telegram/save-image", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
	})
}
