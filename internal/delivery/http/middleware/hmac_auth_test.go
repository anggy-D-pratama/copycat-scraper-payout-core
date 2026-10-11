package middleware

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

func generateSignature(secretKey, apiKey, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(fmt.Sprintf("%s:%s", apiKey, timestamp)))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestHMACAuth_Unit(t *testing.T) {
	apiKey := "my-api-key"
	secretKey := "my-secret-key"
	timestamp := "1773223200000"

	app := fiber.New()
	app.Use(NewHMACAuth(HMACAuthConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	t.Run("Valid Authentication", func(t *testing.T) {
		sig := generateSignature(secretKey, apiKey, timestamp)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Missing x-api-key Header", func(t *testing.T) {
		sig := generateSignature(secretKey, apiKey, timestamp)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Missing x-timestamp Header", func(t *testing.T) {
		sig := generateSignature(secretKey, apiKey, timestamp)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Missing x-signature Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid API Key", func(t *testing.T) {
		sig := generateSignature(secretKey, "wrong-key", timestamp)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", "wrong-key")
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid Signature Mismatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", "deadbeef1234567890abcdef")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid Signature Non-Hex", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", "not-a-hex-string")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Environment Variables Fallback", func(t *testing.T) {
		os.Setenv("API_KEY", "env-key")
		os.Setenv("SECRET_KEY", "env-secret")
		defer func() {
			os.Unsetenv("API_KEY")
			os.Unsetenv("SECRET_KEY")
		}()

		appEnv := fiber.New()
		appEnv.Use(NewHMACAuth())
		appEnv.Get("/env-test", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
		})

		sig := generateSignature("env-secret", "env-key", timestamp)
		req := httptest.NewRequest(http.MethodGet, "/env-test", nil)
		req.Header.Set("x-api-key", "env-key")
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", sig)

		resp, err := appEnv.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
