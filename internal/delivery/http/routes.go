package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, userHandler *UserHandler, withdrawHandler *WithdrawHandler) {
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
	api.Get("/users/:id", userHandler.GetUser)

	// Withdraw Route
	app.Get("/withdraw", withdrawHandler.GetWithdraws)
}
