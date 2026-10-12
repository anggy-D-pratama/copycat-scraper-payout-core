package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type TelegramHandler struct{}

func NewTelegramHandler() *TelegramHandler {
	return &TelegramHandler{}
}

// SendMessage godoc
// @Summary      Mock Telegram Send Message
// @Description  Mock endpoint to receive telegram send requests from bot
// @Tags         Telegram
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /telegram/send [post]
func (h *TelegramHandler) SendMessage(c *fiber.Ctx) error {
	log.Printf("[INFO] Telegram mock SendMessage received request body: %s", string(c.Body()))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

// SaveImage godoc
// @Summary      Mock Telegram Save Image
// @Description  Mock endpoint to receive telegram image/screenshot uploads from bot
// @Tags         Telegram
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Security     Timestamp
// @Security     Signature
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /telegram/save-image [post]
func (h *TelegramHandler) SaveImage(c *fiber.Ctx) error {
	accountNumber := c.FormValue("account_number")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("[WARN] Telegram mock SaveImage received request without file or invalid multipart: %v", err)
	} else {
		log.Printf("[INFO] Telegram mock SaveImage received file: %s (size: %d) for account: %s", fileHeader.Filename, fileHeader.Size, accountNumber)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
