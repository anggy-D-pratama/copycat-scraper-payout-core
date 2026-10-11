package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"boilerplate/internal/model"
	"boilerplate/internal/usecase"
)

type WithdrawHandler struct {
	usecase usecase.WithdrawUsecase
}

func NewWithdrawHandler(u usecase.WithdrawUsecase) *WithdrawHandler {
	return &WithdrawHandler{usecase: u}
}

// GetWithdraws godoc
// @Summary      Get Withdraws
// @Description  Get withdraw transactions with optional filter
// @Tags         Withdraw
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Param        bank_account_id  query     int     false  "Bank Account ID"
// @Param        status           query     string  false  "Status filter"
// @Success      200              {object}  map[string]interface{}
// @Failure      401              {object}  map[string]interface{}
// @Failure      500              {object}  map[string]interface{}
// @Router       /withdraw [get]
func (h *WithdrawHandler) GetWithdraws(c *fiber.Ctx) error {
	bankAccountID := c.QueryInt("bank_account_id", 0)
	status := c.Query("status")

	data, err := h.usecase.GetWithdraws(bankAccountID, status)
	if err != nil {
		log.Printf("[ERROR] GetWithdraws: failed to get withdraws: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}

// UpdateWithdraw godoc
// @Summary      Update Withdraw Status
// @Description  Update withdraw status, notes, or unique code
// @Tags         Withdraw
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Param        body  body      model.WithdrawUpdateRequest  true  "Withdraw Update Payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /withdraw/update [post]
func (h *WithdrawHandler) UpdateWithdraw(c *fiber.Ctx) error {
	var req model.WithdrawUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] UpdateWithdraw: invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.ID == 0 || req.Status == "" {
		log.Printf("[ERROR] UpdateWithdraw: id and status are required, got id=%d, status=%s", req.ID, req.Status)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id and status are required"})
	}

	if err := h.usecase.UpdateWithdraw(req); err != nil {
		log.Printf("[ERROR] UpdateWithdraw: failed to update withdraw id=%d: %v", req.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "withdraw updated"})
}

// ValidateWithdraw godoc
// @Summary      Validate Withdraw
// @Description  Validate withdraw transaction
// @Tags         Withdraw
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Param        body  body      model.WithdrawValidateRequest  true  "Withdraw Validate Payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /withdraw/validate [post]
func (h *WithdrawHandler) ValidateWithdraw(c *fiber.Ctx) error {
	var req model.WithdrawValidateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] ValidateWithdraw: invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.ID == 0 || req.TransactionNo == "" || req.TransactionCode == "" {
		log.Printf("[ERROR] ValidateWithdraw: missing required fields id=%d, transaction_no=%s, transaction_code=%s", req.ID, req.TransactionNo, req.TransactionCode)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id, transaction_no, and transaction_code are required"})
	}

	if err := h.usecase.ValidateWithdraw(req); err != nil {
		log.Printf("[ERROR] ValidateWithdraw: validation failed for id=%d: %v", req.ID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "withdraw validated"})
}
