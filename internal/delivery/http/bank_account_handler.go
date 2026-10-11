package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"boilerplate/internal/usecase"
)

type BankAccountHandler struct {
	usecase usecase.BankAccountUsecase
}

func NewBankAccountHandler(u usecase.BankAccountUsecase) *BankAccountHandler {
	return &BankAccountHandler{usecase: u}
}

// GetBankAccount godoc
// @Summary      Get Bank Account
// @Description  Get bank account detail by account number
// @Tags         Bank Account
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Param        account_number  query     string  true  "Account Number"
// @Success      200             {object}  model.BankAccountResponse
// @Failure      400             {object}  map[string]interface{}
// @Failure      401             {object}  map[string]interface{}
// @Failure      404             {object}  map[string]interface{}
// @Router       /bank-account [get]
func (h *BankAccountHandler) GetBankAccount(c *fiber.Ctx) error {
	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		log.Printf("[ERROR] GetBankAccount: account_number is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "account_number is required"})
	}

	acc, err := h.usecase.GetByAccountNumber(accountNumber)
	if err != nil {
		log.Printf("[ERROR] GetBankAccount: account %s not found: %v", accountNumber, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "bank account not found"})
	}

	return c.JSON(acc)
}

// GetDailyLimit godoc
// @Summary      Get Daily Limit
// @Description  Get daily limit for bank account
// @Tags         Bank Account
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Param        account_number  query     string  true  "Account Number"
// @Success      200             {object}  model.DailyLimitResponse
// @Failure      400             {object}  map[string]interface{}
// @Failure      401             {object}  map[string]interface{}
// @Failure      404             {object}  map[string]interface{}
// @Router       /bank-account/daily-limit [get]
func (h *BankAccountHandler) GetDailyLimit(c *fiber.Ctx) error {
	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		log.Printf("[ERROR] GetDailyLimit: account_number is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "account_number is required"})
	}

	limit, err := h.usecase.GetDailyLimit(accountNumber)
	if err != nil {
		log.Printf("[ERROR] GetDailyLimit: account %s not found: %v", accountNumber, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "bank account not found"})
	}

	return c.JSON(limit)
}
