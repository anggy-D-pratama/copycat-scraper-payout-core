package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

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
	dummyHandler := func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	}

	handlers := &Handlers{
		User:        &UserHandler{},
		Withdraw:    &WithdrawHandler{},
		BankAccount: &BankAccountHandler{},
	}

	// Buat router dengan handlers mock/override untuk route test
	SetupRoutes(app, handlers)

	t.Run("Public Route Health Check - No Auth Header Required", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
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

		req := httptest.NewRequest(http.MethodGet, "/bank-account", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", ts)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		// Request lolos ke handler bankAccount (yang mungkin 500 karena DB belum konek di unit test, tapi bukan 401 Unauthorized)
		assert.NotEqual(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
	_ = dummyHandler
}
