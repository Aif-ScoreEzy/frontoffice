package config

import (
	"front-office/middleware"
	"front-office/pkg/auth"
	"front-office/pkg/company"
	"front-office/pkg/permission"
	"front-office/pkg/product"
	"front-office/pkg/role"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/fo")

	// auth
	api.Post("/register-admin", middleware.IsRequestValid(auth.RegisterAdminRequest{}), auth.RegisterAdmin)
	api.Post("/send-email-verification", middleware.IsRequestValid(auth.SendEmailVerificationRequest{}), auth.SendEmailVerification)
	api.Put("/verify/:token", middleware.SetHeaderAuth, middleware.Auth(), middleware.GetUserIDFromJWT(), auth.VerifyUser)
	api.Post("/request-password-reset", middleware.IsRequestValid(auth.RequestPasswordResetRequest{}), auth.RequestPasswordReset)
	api.Put("/password-reset/:token", middleware.SetHeaderAuth, middleware.Auth(), middleware.GetUserIDFromJWT(), middleware.IsRequestValid(auth.PasswordResetRequest{}), auth.PasswordReset)
	api.Post("/login", middleware.IsRequestValid(auth.UserLoginRequest{}), auth.Login)
	api.Put("/change-password", middleware.Auth(), middleware.IsRequestValid(auth.ChangePasswordRequest{}), middleware.GetUserIDFromJWT(), auth.ChangePassword)

	// user
	api.Post("/register-member", middleware.Auth(), middleware.GetUserIDFromJWT(), middleware.IsRequestValid(user.RegisterMemberRequest{}), user.RegisterMember)
	api.Put("/user/:id", middleware.Auth(), middleware.IsRequestValid(user.UpdateUserRequest{}), user.UpdateUserByID)
	api.Put("/activate/:key", middleware.Auth(), user.ActivateUser)
	api.Put("/deactivate/:email", middleware.Auth(), user.DeactiveUser)
	api.Get("/users", middleware.Auth(), user.GetAllUsers)
	api.Get("/user/:id", middleware.Auth(), user.GetUserByID)

	// company
	api.Put("/company/:id", middleware.Auth(), middleware.IsRequestValid(company.UpdateCompanyRequest{}), company.UpdateCompanyByID)

	// role
	api.Post("/role", middleware.Auth(), middleware.IsRequestValid(role.RoleRequest{}), role.CreateRole)
	api.Get("/roles", middleware.Auth(), role.GetAllRoles)
	api.Get("/role/:id", middleware.Auth(), role.GetRoleByID)
	api.Put("/role/:id", middleware.Auth(), middleware.IsRequestValid(role.RoleRequest{}), role.UpdateRole)
	api.Delete("/role/:id", middleware.Auth(), role.DeleteRole)

	// permission
	api.Post("/permission", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.CreatePermission)
	api.Get("/permission/:id", middleware.Auth(), permission.GetPermissionByID)
	api.Put("/permission/:id", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.UpdatePermissionByID)
	api.Delete("/permission/:id", middleware.Auth(), permission.DeletePermissionByID)

	// product
	api.Post("/product", middleware.Auth(), middleware.IsRequestValid(product.ProductRequest{}), product.CreateProduct)
	api.Get("/products", middleware.Auth(), product.GetAllProducts)
	api.Get("/product/:id", middleware.Auth(), product.GetProductByID)
	api.Put("/product/:id", middleware.Auth(), middleware.IsRequestValid(product.UpdateProductRequest{}), product.UpdateProductByID)
	api.Delete("/product/:id", middleware.Auth(), product.DeleteProductByID)
}
