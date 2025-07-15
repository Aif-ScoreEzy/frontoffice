package grading

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(gradingAPI fiber.Router, cfg *config.Config) {
	repo := NewRepository(cfg)
	service := NewService(repo)
	controller := NewController(service)

	gradingAPI.Get("/", middleware.GetJWTPayloadFromCookie(), controller.GetGradings)
}
