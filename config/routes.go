package config

import (
	"front-office/middleware"
	"front-office/pkg/auth"
	"front-office/pkg/company"
	genRetail "front-office/pkg/gen-retail"
	"front-office/pkg/grading"
	"front-office/pkg/log"
	"front-office/pkg/permission"
	"front-office/pkg/role"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/fo")

	// auth
	api.Post("/register-admin", middleware.IsRequestValid(auth.RegisterAdminRequest{}), auth.RegisterAdmin)
	api.Put("/verify/:token", middleware.SetHeaderAuth, middleware.IsRequestValid(auth.PasswordResetRequest{}), auth.VerifyUser)
	api.Post("/request-password-reset", middleware.IsRequestValid(auth.RequestPasswordResetRequest{}), auth.RequestPasswordReset)
	api.Put("/password-reset/:token", middleware.SetHeaderAuth, middleware.GetPayloadFromJWT(), middleware.IsRequestValid(auth.PasswordResetRequest{}), auth.PasswordReset)
	api.Post("/login", middleware.IsRequestValid(auth.UserLoginRequest{}), auth.Login)
	api.Put("/change-password", middleware.Auth(), middleware.IsRequestValid(auth.ChangePasswordRequest{}), middleware.GetPayloadFromJWT(), auth.ChangePassword)

	// user
	api.Post("/register-member", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), middleware.IsRequestValid(user.RegisterMemberRequest{}), auth.RegisterMember)
	api.Get("/users", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), user.GetAllUsers)
	api.Get("/user/:id", middleware.Auth(), middleware.GetPayloadFromJWT(), user.GetUserByID)
	api.Put("/send-email-activation/:email", middleware.Auth(), middleware.AdminAuth(), middleware.GetPayloadFromJWT(), auth.SendEmailActivation)
	api.Put("/user/:id", middleware.AdminAuth(), middleware.IsRequestValid(user.UpdateUserRequest{}), middleware.GetPayloadFromJWT(), user.UpdateUserByID)
	api.Put("/edit-profile", middleware.Auth(), middleware.IsRequestValid(user.UpdateProfileRequest{}), middleware.GetPayloadFromJWT(), user.UpdateProfile)
	api.Put("/upload-profile-image", middleware.Auth(), middleware.IsRequestValid(user.UploadProfileImageRequest{}), middleware.GetPayloadFromJWT(), middleware.FileUpload(), user.UploadProfileImage)
	api.Delete("/user/:id", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), user.DeleteUserByID)

	// grading
	api.Post("/create-gradings", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), middleware.IsRequestValid(grading.CreateGradingsRequest{}), grading.CreateGradings)
	api.Get("/get-gradings", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), grading.GetGradings)
	api.Put("/update-gradings", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), middleware.IsRequestValid(grading.CreateGradingsRequest{}), grading.ReplaceGradings)
	api.Put("/update-gradings-new", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), middleware.IsRequestValid(grading.CreateGradingsNewRequest{}), grading.ReplaceGradingsNew)
	// score
	api.Post("/request-score", middleware.GetPayloadFromJWT(), middleware.IsRequestValid(genRetail.GenRetailRequest{}), genRetail.RequestScore)
	api.Get("/download-csv/:opsi", middleware.GetPayloadFromJWT(), genRetail.DownloadCSV)

	// company
	api.Put("/company/:id", middleware.Auth(), middleware.IsRequestValid(company.UpdateCompanyRequest{}), company.UpdateCompanyByID)

	// role
	api.Post("/role", middleware.Auth(), middleware.IsRequestValid(role.CreateRoleRequest{}), role.CreateRole)
	api.Get("/roles", middleware.Auth(), role.GetAllRoles)
	api.Get("/role/:id", middleware.Auth(), role.GetRoleByID)
	api.Put("/role/:id", middleware.Auth(), middleware.IsRequestValid(role.UpdateRoleRequest{}), role.UpdateRole)
	api.Delete("/role/:id", middleware.Auth(), role.DeleteRole)

	// permission
	api.Post("/permission", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.CreatePermission)
	api.Get("/permission/:id", middleware.Auth(), permission.GetPermissionByID)
	api.Put("/permission/:id", middleware.Auth(), middleware.IsRequestValid(permission.PermissionRequest{}), permission.UpdatePermissionByID)
	api.Delete("/permission/:id", middleware.Auth(), permission.DeletePermissionByID)

	// log
	api.Get("/get-logs-by-date", log.GetTransactionLogsByDate)
	api.Get("/get-logs-by-range-date", log.GetTransactionLogsByRangeDate)

	// product
	// api.Post("/product", middleware.Auth(), middleware.IsRequestValid(product.ProductRequest{}), product.CreateProduct)
	// api.Get("/products", middleware.Auth(), product.GetAllProducts)
	// api.Get("/product/:id", middleware.Auth(), product.GetProductByID)
	// api.Put("/product/:id", middleware.Auth(), middleware.IsRequestValid(product.UpdateProductRequest{}), product.UpdateProductByID)
	// api.Delete("/product/:id", middleware.Auth(), product.DeleteProductByID)

}
