package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// HMACAuthConfig mendefinisikan konfigurasi untuk HMACAuth middleware.
type HMACAuthConfig struct {
	APIKey    string
	SecretKey string
}

// NewHMACAuth membuat middleware autentikasi HMAC-SHA256 berdasarkan konfigurasi yang diberikan.
// Jika config tidak dispesifikasikan (atau string kosong), nilai dibaca dari env API_KEY dan SECRET_KEY.
func NewHMACAuth(config ...HMACAuthConfig) fiber.Handler {
	var cfg HMACAuthConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("API_KEY")
	}
	if cfg.SecretKey == "" {
		cfg.SecretKey = os.Getenv("SECRET_KEY")
	}

	return func(c *fiber.Ctx) error {
		apiKeyHeader := c.Get("x-api-key")
		timestampHeader := c.Get("x-timestamp")
		signatureHeader := c.Get("x-signature")

		// 1. Validasi keberadaan ketiga header
		if apiKeyHeader == "" || timestampHeader == "" || signatureHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Unauthorized: missing required authentication headers (x-api-key, x-timestamp, x-signature)",
				},
			})
		}

		// 2. Validasi x-api-key
		if apiKeyHeader != cfg.APIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Unauthorized: invalid api key",
				},
			})
		}

		// 3. Hitung expected signature: HMAC-SHA256(SECRET_KEY, `${API_KEY}:${timestamp}`)
		mac := hmac.New(sha256.New, []byte(cfg.SecretKey))
		payload := fmt.Sprintf("%s:%s", apiKeyHeader, timestampHeader)
		mac.Write([]byte(payload))
		expectedMAC := mac.Sum(nil)
		expectedHex := hex.EncodeToString(expectedMAC)

		// 4. Bandingkan x-signature dengan expected signature menggunakan constant time comparison
		providedMAC, err := hex.DecodeString(signatureHeader)
		if err != nil || !hmac.Equal(providedMAC, expectedMAC) || signatureHeader != expectedHex {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Unauthorized: invalid signature",
				},
			})
		}

		return c.Next()
	}
}
