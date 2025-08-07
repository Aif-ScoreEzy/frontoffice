package grade

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(gradingAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	service := NewService(repo)
	controller := NewController(service)

	gradingAPI.Put("/", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(createGradeRequest{}), controller.SaveGrading)
	gradingAPI.Get("/", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetGrades)
}
