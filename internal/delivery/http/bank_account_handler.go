package http

import (
	"github.com/gofiber/fiber/v2"
	"boilerplate/internal/service"
)

type BankAccountHandler struct {
	service service.BankAccountService
}

func NewBankAccountHandler(s service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{service: s}
}

func (h *BankAccountHandler) GetBankAccount(c *fiber.Ctx) error {
	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "account_number is required"})
	}

	acc, err := h.service.GetByAccountNumber(accountNumber)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "bank account not found"})
	}

	return c.JSON(acc)
}

func (h *BankAccountHandler) GetDailyLimit(c *fiber.Ctx) error {
	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "account_number is required"})
	}

	limit, err := h.service.GetDailyLimit(accountNumber)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "bank account not found"})
	}

	return c.JSON(limit)
}
