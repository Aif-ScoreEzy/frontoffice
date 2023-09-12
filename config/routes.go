package config

import (
	"front-office/middleware"
	"front-office/pkg/company"
	"front-office/pkg/permission"
	"front-office/pkg/product"
	"front-office/pkg/role"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/fo")

	api.Post("/role", middleware.Auth(), middleware.IsRequestValid(role.RoleRequest{}), role.CreateRole)
	api.Get("/roles", middleware.Auth(), role.GetAllRoles)
	api.Get("/role/:id", middleware.Auth(), role.GetRoleByID)
	api.Put("/role/:id", middleware.Auth(), middleware.IsRequestValid(role.RoleRequest{}), role.UpdateRole)
	api.Delete("/role/:id", middleware.Auth(), role.DeleteRole)

	api.Post("/permission", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.CreatePermission)
	api.Get("/permission/:id", middleware.Auth(), permission.GetPermissionByID)
	api.Put("/permission/:id", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.UpdatePermissionByID)
	api.Delete("/permission/:id", middleware.Auth(), permission.DeletePermissionByID)

	api.Put("/company/:id", middleware.Auth(), middleware.IsRequestValid(company.UpdateCompanyRequest{}), company.UpdateCompanyByID)

	api.Post("/register", middleware.IsRequestValid(user.RegisterUserRequest{}), user.Register)
	api.Post("/register-member", middleware.Auth(), middleware.GetUserIDFromJWT(), middleware.IsRequestValid(user.RegisterMemberRequest{}), user.RegisterMember)
	api.Post("/send-email-verification", middleware.IsRequestValid(user.SendEmailVerificationRequest{}), user.SendEmailVerification)
	api.Put("/verify/:token", middleware.SetHeaderAuth, middleware.Auth(), middleware.GetUserIDFromJWT(), user.VerifyUser)
	api.Post("/request-password-reset", middleware.IsRequestValid(user.RequestPasswordResetRequest{}), user.RequestPasswordReset)
	api.Put("/password-reset/:token", middleware.SetHeaderAuth, middleware.Auth(), middleware.GetUserIDFromJWT(), middleware.IsRequestValid(user.PasswordResetRequest{}), user.PasswordReset)
	api.Post("/login", middleware.IsRequestValid(user.UserLoginRequest{}), user.Login)
	api.Put("/user/:id", middleware.Auth(), middleware.IsRequestValid(user.UpdateUserRequest{}), user.UpdateUserByID)
	api.Put("/activate/:key", middleware.Auth(), user.ActivateUser)
	api.Put("/deactivate/:email", middleware.Auth(), user.DeactiveUser)
	api.Get("/users", middleware.Auth(), user.GetAllUsers)
	api.Get("/user/:id", middleware.Auth(), user.GetUserByID)

	api.Post("/product", middleware.Auth(), middleware.IsRequestValid(product.ProductRequest{}), product.CreateProduct)
	api.Get("/products", middleware.Auth(), product.GetAllProducts)
	api.Get("/product/:id", middleware.Auth(), product.GetProductByID)
	api.Put("/product/:id", middleware.Auth(), middleware.IsRequestValid(product.UpdateProductRequest{}), product.UpdateProductByID)
	api.Delete("/product/:id", middleware.Auth(), product.DeleteProductByID)
}
