package config

import (
	"front-office/middleware"
	"front-office/pkg/permission"
	"front-office/pkg/role"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/role", middleware.IsRoleRequestValid, role.CreateRole)
	api.Get("/roles", role.GetAllRoles)
	api.Get("/role/:id", role.GetRoleByID)
	api.Put("/role/:id", middleware.IsRoleRequestValid, role.UpdateRole)
	api.Delete("/role/:id", role.DeleteRole)

	api.Post("/permission", middleware.IsPermissionRequestValid, permission.CreatePermission)
	api.Get("/permission/:id", permission.GetRoleByID)
	api.Put("/permission/:id", middleware.IsPermissionRequestValid, permission.UpdatePermissionByID)
	api.Delete("/permission/:id", permission.DeletePermissionByID)

	api.Post("/register", middleware.IsRegisterUserRequestValid, user.Register)
}
