package auth

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/member"
	"front-office/pkg/core/passwordresettoken"
	"front-office/pkg/core/role"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(authAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	memberRepo := member.NewRepository(cfg, client, nil)
	roleRepo := role.NewRepository(cfg, client)
	activationTokenRepo := activationtoken.NewRepository(cfg, client, nil)
	passwordResetRepo := passwordresettoken.NewRepository(cfg, client, nil)
	logOperationRepo := operation.NewRepository(cfg, client, nil)

	service := NewService(cfg, repo, memberRepo, roleRepo, logOperationRepo, activationTokenRepo, passwordResetRepo)
	serviceUser := member.NewService(memberRepo, roleRepo, logOperationRepo)
	serviceActivationToken := activationtoken.NewService(activationTokenRepo, cfg)
	servicePasswordResetToken := passwordresettoken.NewService(passwordResetRepo, cfg)
	serviceLogOperation := operation.NewService(logOperationRepo)

	controller := NewController(service, serviceUser, serviceActivationToken, servicePasswordResetToken, serviceLogOperation, cfg)

	authAPI.Post("/register-member", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(member.RegisterMemberRequest{}), controller.RegisterMember)
	authAPI.Post("/login", middleware.IsRequestValid(userLoginRequest{}), controller.Login)
	authAPI.Put("/verify/:token", middleware.SetHeaderAuth, middleware.IsRequestValid(PasswordResetRequest{}), controller.VerifyUser)
	authAPI.Post("/logout", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.Logout)
	authAPI.Post("/refresh-access", middleware.GetPayloadFromRefreshToken(), controller.RefreshAccessToken)
	authAPI.Put("/send-email-activation/:email", middleware.Auth(), middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), controller.RequestActivation)
	authAPI.Post("/request-password-reset", middleware.IsRequestValid(RequestPasswordResetRequest{}), controller.RequestPasswordReset)
	authAPI.Put("/password-reset/:token", middleware.SetCookiePasswordResetToken, middleware.GetJWTPayloadPasswordResetFromCookie(), middleware.IsRequestValid(PasswordResetRequest{}), controller.PasswordReset)
	authAPI.Put("/change-password", middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(ChangePasswordRequest{}), controller.ChangePassword)
}
