package app

import (
	"log"
	"os"
	"strings"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	httpDelivery "boilerplate/internal/delivery/http"
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
	"boilerplate/internal/usecase"
	"boilerplate/pkg/database"
)

func Run() {
	// Init DB
	db := database.Connect()
	if err := db.AutoMigrate(&model.BankAccount{}, &model.Withdraw{}, &model.BankTransaction{}); err != nil {
		log.Fatal("AutoMigrate failed: ", err)
	}

	// Seed Dummy Data
	seed(db)

	server := fiber.New()

	// Dependency Injection
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := httpDelivery.NewUserHandler(userUsecase)

	bankAccountRepo := repository.NewBankAccountRepository(db)
	bankAccountUsecase := usecase.NewBankAccountUsecase(bankAccountRepo)
	bankAccountHandler := httpDelivery.NewBankAccountHandler(bankAccountUsecase)

	withdrawRepo := repository.NewWithdrawRepository(db)
	withdrawUsecase := usecase.NewWithdrawUsecase(withdrawRepo)
	withdrawHandler := httpDelivery.NewWithdrawHandler(withdrawUsecase)

	bankTransactionRepo := repository.NewBankTransactionRepository(db)
	bankTransactionUsecase := usecase.NewBankTransactionUsecase(bankTransactionRepo)
	bankTransactionHandler := httpDelivery.NewBankTransactionHandler(bankTransactionUsecase)

	telegramHandler := httpDelivery.NewTelegramHandler()

	handlers := &httpDelivery.Handlers{
		User:            userHandler,
		Withdraw:        withdrawHandler,
		BankAccount:     bankAccountHandler,
		BankTransaction: bankTransactionHandler,
		Telegram:        telegramHandler,
	}

	httpDelivery.SetupRoutes(server, handlers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Fatal(server.Listen(port))
}

func seed(db *gorm.DB) {
	// Basic seeder for testing bot
	var count int64
	db.Model(&model.BankAccount{}).Count(&count)
	if count == 0 {
		acc := model.BankAccount{
			AccountNumber: "1234567890",
			BankName:      "BCA",
			AccountName:   "Bot Testing Account",
			DeviceID:      "device-123",
			DeviceCode:    "device-123",
			PortAppium:    4723,
			SafeKey:       "123456",
			Pin:           "123456",
			Username:      "bot_user",
			Password:      "password",
			DailyLimit:    100000000,
		}
		db.Create(&acc)

		db.Create(&model.Withdraw{
			TransactionNo:   "TRX-001",
			TransactionCode: "TXC-001",
			BankCode:        "BCA",
			Amount:          50000,
			Remark:          "Test withdraw 1",
			BankAccount:     "1234567890",
			BankAccountId:   acc.ID,
			Status:          "queue",
		})
		db.Create(&model.Withdraw{
			TransactionNo:   "TRX-002",
			TransactionCode: "TXC-002",
			BankCode:        "BCA",
			Amount:          150000,
			Remark:          "Test withdraw 2",
			BankAccount:     "1234567890",
			BankAccountId:   acc.ID,
			Status:          "queue",
		})
	}
}
