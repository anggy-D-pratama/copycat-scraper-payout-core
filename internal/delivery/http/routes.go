package http

import (
	"boilerplate/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handlers struct {
	User        *UserHandler
	Withdraw    *WithdrawHandler
	BankAccount *BankAccountHandler
}

func SetupRoutes(app *fiber.App, h *Handlers) {
	// Default Health Check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Service is healthy",
		})
	})

	// Default Swagger Docs
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API Routes
	api := app.Group("/api/v1")
	api.Get("/users/:id", h.User.GetUser)

	// HMAC Authentication Middleware untuk Endpoint Core
	hmacAuth := middleware.NewHMACAuth()

	// Bank Account Routes (Protected)
	bankAccount := app.Group("/bank-account", hmacAuth)
	bankAccount.Get("", h.BankAccount.GetBankAccount)
	bankAccount.Get("/daily-limit", h.BankAccount.GetDailyLimit)

	// Withdraw Routes (Protected)
	withdraw := app.Group("/withdraw", hmacAuth)
	withdraw.Get("", h.Withdraw.GetWithdraws)
	withdraw.Post("/update", h.Withdraw.UpdateWithdraw)
	withdraw.Post("/validate", h.Withdraw.ValidateWithdraw)
}
