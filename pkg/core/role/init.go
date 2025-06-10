package role

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(roleAPI fiber.Router, cfg *config.Config) {
	repo := NewRepository(cfg)
	service := NewService(repo)
	controller := NewController(service)

	roleAPI.Get("/", middleware.Auth(), controller.GetAllRoles)
	roleAPI.Get("/:id", middleware.Auth(), controller.GetRoleById)
}
