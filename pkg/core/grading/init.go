package grading

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(gradingAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	service := NewService(repo)
	controller := NewController(service)

	gradingAPI.Post("/", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(CreateGradingsRequest{}), controller.CreateGradings)
	gradingAPI.Get("/", middleware.GetJWTPayloadFromCookie(), controller.GetGradings)
	gradingAPI.Put("/", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(CreateGradingsNewRequest{}), controller.ReplaceGradingsNew)
}
