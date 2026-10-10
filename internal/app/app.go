package app

import (
	"log"

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
	if err := db.AutoMigrate(&model.Withdraw{}); err != nil {
		log.Fatal("AutoMigrate failed: ", err)
	}

	server := fiber.New()

	// Dependency Injection
	userRepo := repository.NewUserRepository()
	userSvc := service.NewUserService(userRepo)
	userHandler := httpDelivery.NewUserHandler(userSvc)

	withdrawRepo := repository.NewWithdrawRepository(db)
	withdrawSvc := service.NewWithdrawService(withdrawRepo)
	withdrawHandler := httpDelivery.NewWithdrawHandler(withdrawSvc)

	httpDelivery.SetupRoutes(server, userHandler, withdrawHandler)

	log.Fatal(server.Listen(":3000"))
}
