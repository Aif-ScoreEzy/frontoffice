package config

import (
	"front-office/pkg/role"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/role", role.CreateRole)
}
