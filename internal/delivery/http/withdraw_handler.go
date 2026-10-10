package http

import (
	"github.com/gofiber/fiber/v2"
	"boilerplate/internal/model"
	"boilerplate/internal/service"
)

type WithdrawHandler struct {
	service service.WithdrawService
}

func NewWithdrawHandler(s service.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{service: s}
}

func (h *WithdrawHandler) GetWithdraws(c *fiber.Ctx) error {
	bankAccountID := c.QueryInt("bank_account_id", 0)
	status := c.Query("status")

	data, err := h.service.GetWithdraws(bankAccountID, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}

func (h *WithdrawHandler) UpdateWithdraw(c *fiber.Ctx) error {
	var req model.WithdrawUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.service.UpdateWithdraw(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "withdraw updated"})
}

func (h *WithdrawHandler) ValidateWithdraw(c *fiber.Ctx) error {
	var req model.WithdrawValidateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.service.ValidateWithdraw(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "withdraw validated"})
}
