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

	api.Post("/role", middleware.Auth(), middleware.IsRoleRequestValid, role.CreateRole)
	api.Get("/roles", middleware.Auth(), role.GetAllRoles)
	api.Get("/role/:id", middleware.Auth(), role.GetRoleByID)
	api.Put("/role/:id", middleware.Auth(), middleware.IsRoleRequestValid, role.UpdateRole)
	api.Delete("/role/:id", middleware.Auth(), role.DeleteRole)

	api.Post("/permission", middleware.Auth(), middleware.IsPermissionRequestValid, permission.CreatePermission)
	api.Get("/permission/:id", middleware.Auth(), permission.GetRoleByID)
	api.Put("/permission/:id", middleware.Auth(), middleware.IsPermissionRequestValid, permission.UpdatePermissionByID)
	api.Delete("/permission/:id", middleware.Auth(), permission.DeletePermissionByID)

	api.Post("/register", middleware.IsRegisterUserRequestValid, user.Register)
	api.Post("/login", middleware.IsLoginRequestValid, user.Login)
	api.Put("/user/:id", middleware.Auth(), middleware.IsUpdateUserRequestValid, user.UpdateUserByID)
}
