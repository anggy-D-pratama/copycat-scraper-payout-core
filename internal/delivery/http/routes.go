package http

import (
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

	// Bank Account Routes
	app.Get("/bank-account", h.BankAccount.GetBankAccount)
	app.Get("/bank-account/daily-limit", h.BankAccount.GetDailyLimit)

	// Withdraw Routes
	app.Get("/withdraw", h.Withdraw.GetWithdraws)
	app.Post("/withdraw/update", h.Withdraw.UpdateWithdraw)
	app.Post("/withdraw/validate", h.Withdraw.ValidateWithdraw)
}
