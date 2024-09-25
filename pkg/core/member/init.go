package member

import (
	"front-office/app/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(userAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	service := NewService(repo)
	controller := NewController(service)

	// userAPI.Get("/:id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetBy)
	userAPI.Get("/", controller.GetBy)
}
