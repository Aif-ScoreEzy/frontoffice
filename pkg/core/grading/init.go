package grading

import (
	"front-office/app/config"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(gradingAPI fiber.Router, cfg *config.Config) {
	// repo := NewRepository(cfg)
	// service := NewService(repo)
	// controller := NewController(service)

	// gradingAPI.Post("/", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(CreateGradingsRequest{}), controller.CreateGradings)
	// gradingAPI.Get("/", middleware.GetJWTPayloadFromCookie(), controller.GetGradings)
	// gradingAPI.Put("/", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(CreateGradingsNewRequest{}), controller.ReplaceGradingsNew)
}
