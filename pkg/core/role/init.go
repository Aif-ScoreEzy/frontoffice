package role

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(roleAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client)
	service := NewService(repo)
	controller := NewController(service)

	roleAPI.Get("/", middleware.Auth(), controller.GetRoles)
	roleAPI.Get("/:id", middleware.Auth(), controller.GetRoleById)
}
