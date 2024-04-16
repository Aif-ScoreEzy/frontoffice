package grading

import (
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(gradingAPI fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewController(service)

	gradingAPI.Post("/", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), middleware.IsRequestValid(CreateGradingsRequest{}), controller.CreateGradings)
	gradingAPI.Get("/", middleware.GetPayloadFromJWT(), controller.GetGradings)
	gradingAPI.Put("/", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), middleware.IsRequestValid(CreateGradingsNewRequest{}), controller.ReplaceGradingsNew)
}
