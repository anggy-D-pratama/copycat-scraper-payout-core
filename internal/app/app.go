package app

import (
	"log"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	httpDelivery "boilerplate/internal/delivery/http"
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
	"boilerplate/internal/service"
	"boilerplate/pkg/database"
)

func Run() {
	// Init DB
	db := database.Connect()
	if err := db.AutoMigrate(&model.BankAccount{}, &model.Withdraw{}); err != nil {
		log.Fatal("AutoMigrate failed: ", err)
	}

	// Seed Dummy Data
	seed(db)

	server := fiber.New()

	// Dependency Injection
	userRepo := repository.NewUserRepository()
	userSvc := service.NewUserService(userRepo)
	userHandler := httpDelivery.NewUserHandler(userSvc)

	bankAccountRepo := repository.NewBankAccountRepository(db)
	bankAccountSvc := service.NewBankAccountService(bankAccountRepo)
	bankAccountHandler := httpDelivery.NewBankAccountHandler(bankAccountSvc)

	withdrawRepo := repository.NewWithdrawRepository(db)
	withdrawSvc := service.NewWithdrawService(withdrawRepo)
	withdrawHandler := httpDelivery.NewWithdrawHandler(withdrawSvc)

	httpDelivery.SetupRoutes(server, userHandler, withdrawHandler, bankAccountHandler)

	log.Fatal(server.Listen(":3000"))
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
			Pin:           "123456",
			Username:      "bot_user",
			Password:      "password",
			DailyLimit:    100000000,
		}
		db.Create(&acc)

		db.Create(&model.Withdraw{
			TransactionNo: "TRX-001",
			Amount:        50000,
			Remark:        "Test withdraw 1",
			BankAccount:   "1234567890",
			BankAccountId: acc.ID,
			Status:        "QUEUE",
		})
		db.Create(&model.Withdraw{
			TransactionNo: "TRX-002",
			Amount:        150000,
			Remark:        "Test withdraw 2",
			BankAccount:   "1234567890",
			BankAccountId: acc.ID,
			Status:        "QUEUE",
		})
	}
}
