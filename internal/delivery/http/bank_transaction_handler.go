package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"boilerplate/internal/model"
	"boilerplate/internal/usecase"
)

type BankTransactionHandler struct {
	usecase usecase.BankTransactionUsecase
}

func NewBankTransactionHandler(u usecase.BankTransactionUsecase) *BankTransactionHandler {
	return &BankTransactionHandler{usecase: u}
}

// BulkInsert godoc
// @Summary      Bulk Bank Transaction Intake
// @Description  Receive scraped bank transaction mutations in bulk
// @Tags         Bank Transaction
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Param        body  body      []model.BankTransactionRequest  true  "Bulk Bank Transaction Payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /core-mobile/bank-transactions/bulk [post]
func (h *BankTransactionHandler) BulkInsert(c *fiber.Ctx) error {
	var requests []model.BankTransactionRequest
	if err := c.BodyParser(&requests); err != nil {
		log.Printf("[ERROR] BulkInsert: invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   "invalid request body: must be a JSON array",
		})
	}

	if len(requests) == 0 {
		log.Printf("[ERROR] BulkInsert: empty request array")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   "payload array tidak boleh kosong",
		})
	}

	res, err := h.usecase.BulkInsert(requests)
	if err != nil {
		log.Printf("[ERROR] BulkInsert: failed to bulk insert: %v", err)
		// Check if validation error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   res,
	})
}
