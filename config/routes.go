package config

import (
	"front-office/middleware"
	"front-office/pkg/permission"
	"front-office/pkg/role"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/role", middleware.IsRoleRequestValid, role.CreateRole)
	api.Get("/role/:id", role.GetRoleByID)
	api.Put("/role/:id", middleware.IsRoleRequestValid, role.UpdateRole)
	api.Delete("/role/:id", role.DeleteRole)

	api.Post("/permission", middleware.IsPermissionRequestValid, permission.CreatePermission)
}
