package auth

import (
	"front-office/app/config"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/member"
	"front-office/pkg/core/passwordresettoken"
	"front-office/pkg/core/role"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(authAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	repoUser := member.NewRepository(db, cfg)
	repoRole := role.NewRepository(db, cfg)
	repoActivationToken := activationtoken.NewRepository(cfg)
	repoPasswordResetToken := passwordresettoken.NewRepository(db, cfg)
	repoLogOperation := operation.NewRepository(cfg)

	service := NewService(repo, repoUser, repoRole, cfg)
	serviceRole := role.NewService(repoRole)
	serviceUser := member.NewService(repoUser, serviceRole)
	serviceActivationToken := activationtoken.NewService(repoActivationToken, cfg)
	servicePasswordResetToken := passwordresettoken.NewService(repoPasswordResetToken, cfg)
	serviceLogOperation := operation.NewService(repoLogOperation)

	controller := NewController(service, serviceUser, serviceActivationToken, servicePasswordResetToken, serviceLogOperation, cfg)

	authAPI.Post("/register-member", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(member.RegisterMemberRequest{}), controller.RegisterMember)
	authAPI.Post("/request-password-reset", middleware.IsRequestValid(RequestPasswordResetRequest{}), controller.RequestPasswordReset)
	authAPI.Post("/login", middleware.IsRequestValid(UserLoginRequest{}), controller.Login)
	authAPI.Post("/logout", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.Logout)
	authAPI.Post("/refresh-access", middleware.GetPayloadFromRefreshToken(), controller.RefreshAccessToken)
	authAPI.Put("/change-password", middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(ChangePasswordRequest{}), controller.ChangePassword)
	authAPI.Put("/send-email-activation/:email", middleware.Auth(), middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), controller.SendEmailActivation)
	authAPI.Put("/verify/:token", middleware.SetHeaderAuth, middleware.IsRequestValid(PasswordResetRequest{}), controller.VerifyUser)
	authAPI.Put("/password-reset/:token", middleware.SetCookiePasswordResetToken, middleware.GetJWTPayloadPasswordResetFromCookie(), middleware.IsRequestValid(PasswordResetRequest{}), controller.PasswordReset)
}
