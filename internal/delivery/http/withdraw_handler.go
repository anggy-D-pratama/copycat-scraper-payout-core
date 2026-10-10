package http

import (
	"github.com/gofiber/fiber/v2"
	"boilerplate/internal/service"
)

type WithdrawHandler struct {
	service service.WithdrawService
}

func NewWithdrawHandler(s service.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{service: s}
}

// GetWithdraws godoc
// @Summary Get list of withdraws
// @Description Retrieve filtered withdraw list by bank_account_id and status
// @Tags withdraw
// @Produce json
// @Param bank_account_id query int false "Bank Account ID"
// @Param status query string false "Withdraw Status"
// @Success 200 {object} map[string][]model.WithdrawResponse
// @Failure 500 {object} map[string]string
// @Router /withdraw [get]
func (h *WithdrawHandler) GetWithdraws(c *fiber.Ctx) error {
	bankAccountID := c.QueryInt("bank_account_id", 0)
	status := c.Query("status")

	data, err := h.service.GetWithdraws(bankAccountID, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}
