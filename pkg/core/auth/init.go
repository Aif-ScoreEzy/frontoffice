package auth

import (
	"front-office/app/config"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/passwordresettoken"
	"front-office/pkg/core/role"
	"front-office/pkg/core/user"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(authAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	repoUser := user.NewRepository(db, cfg)
	repoRole := role.NewRepository(db)
	repoActivationToken := activationtoken.NewRepository(db, cfg)
	repoPasswordResetToken := passwordresettoken.NewRepository(db)

	service := NewService(repo, repoUser, repoRole, cfg)
	serviceUser := user.NewService(repoUser, repoRole)
	serviceActivationToken := activationtoken.NewService(repoActivationToken, cfg)
	servicePasswordResetToken := passwordresettoken.NewService(repoPasswordResetToken, cfg)

	controller := NewController(service, serviceUser, serviceActivationToken, servicePasswordResetToken, cfg)

	authAPI.Post("/register-admin", middleware.IsRequestValid(RegisterAdminRequest{}), controller.RegisterAdmin)
	// authAPI.Post("/register-member", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(user.RegisterMemberRequest{}), controller.RegisterMember)
	authAPI.Post("/register-member", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(user.RegisterMemberRequest{}), controller.RegisterMemberAifCore)
	// authAPI.Post("/request-password-reset", middleware.IsRequestValid(RequestPasswordResetRequest{}), controller.RequestPasswordReset)
	authAPI.Post("/login", middleware.IsRequestValid(UserLoginRequest{}), controller.LoginAifCore)
	authAPI.Post("/logout", controller.Logout)
	// authAPI.Post("/refresh-access", middleware.GetPayloadFromRefreshToken(), controller.RefreshAccessToken)
	// authAPI.Put("/change-password", middleware.IsRequestValid(ChangePasswordRequest{}), controller.ChangePasswordAifcore)
	authAPI.Put("/send-email-activation/:email", middleware.Auth(), middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), controller.SendEmailActivation)
	authAPI.Put("/verify/:token", middleware.SetHeaderAuth, middleware.IsRequestValid(PasswordResetRequest{}), controller.VerifyUser)
	authAPI.Put("/password-reset/:token", middleware.SetHeaderAuth, middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(PasswordResetRequest{}), controller.PasswordReset)
}
