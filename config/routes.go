package config

import (
	"front-office/middleware"
	"front-office/pkg/company"
	"front-office/pkg/permission"
	"front-office/pkg/role"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/role", middleware.Auth(), middleware.IsRequestValid(role.RoleRequest{}), role.CreateRole)
	api.Get("/roles", middleware.Auth(), role.GetAllRoles)
	api.Get("/role/:id", middleware.Auth(), role.GetRoleByID)
	api.Put("/role/:id", middleware.Auth(), middleware.IsRequestValid(role.RoleRequest{}), role.UpdateRole)
	api.Delete("/role/:id", middleware.Auth(), role.DeleteRole)

	api.Post("/permission", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.CreatePermission)
	api.Get("/permission/:id", middleware.Auth(), permission.GetRoleByID)
	api.Put("/permission/:id", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.UpdatePermissionByID)
	api.Delete("/permission/:id", middleware.Auth(), permission.DeletePermissionByID)

	api.Put("company/:id", middleware.Auth(), middleware.IsRequestValid(company.UpdateCompanyRequest{}), company.UpdateCompanyByID)

	api.Post("/register", middleware.IsRequestValid(user.RegisterUserRequest{}), user.Register)
	api.Post("/login", middleware.IsRequestValid(user.UserLoginRequest{}), user.Login)
	api.Put("/user/:id", middleware.Auth(), middleware.IsRequestValid(user.UpdateUserRequest{}), user.UpdateUserByID)
}
