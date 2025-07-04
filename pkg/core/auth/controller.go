package auth

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/member"
	"front-office/pkg/core/passwordresettoken"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func NewController(
	service Service,
	svcUser member.Service,
	svcActivationToken activationtoken.Service,
	svcPasswordResetToken passwordresettoken.Service,
	svcLogOperation operation.Service,
	cfg *config.Config,
) Controller {
	return &controller{
		svc:                   service,
		svcUser:               svcUser,
		svcActivationToken:    svcActivationToken,
		svcPasswordResetToken: svcPasswordResetToken,
		svcLogOperation:       svcLogOperation,
		cfg:                   cfg,
	}
}

type controller struct {
	svc                   Service
	svcUser               member.Service
	svcActivationToken    activationtoken.Service
	svcPasswordResetToken passwordresettoken.Service
	svcLogOperation       operation.Service
	cfg                   *config.Config
}

type Controller interface {
	RegisterMember(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	VerifyUser(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RequestActivation(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	RequestPasswordReset(c *fiber.Ctx) error
	PasswordReset(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
}

func (ctrl *controller) RegisterMember(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*member.RegisterMemberRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	currentUserId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidUserSession)
	}

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidCompanySession)
	}

	reqBody.CompanyId = companyId
	reqBody.RoleId = uint(memberRoleId)

	if err := ctrl.svc.AddMember(currentUserId, reqBody); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", reqBody.Email),
		nil,
	))
}

func (ctrl *controller) VerifyUser(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*PasswordResetRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	token := c.Params("token")
	if token == "" {
		return apperror.BadRequest("missing activation token")
	}

	if err := ctrl.svc.VerifyMember(token, reqBody); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"your account has been verified",
		nil,
	))
}

func (ctrl *controller) Logout(c *fiber.Ctx) error {
	memberId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidUserSession)
	}

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidCompanySession)
	}

	// Clear access & refresh token cookies
	clearAuthCookie(c, "aif_token")
	clearAuthCookie(c, "aif_refresh_token")

	err = ctrl.svc.Logout(memberId, companyId)
	if err != nil {
		log.Warn().Err(err).Msg("failed to log logout event")
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to logout",
		nil,
	))
}

func (ctrl *controller) RequestActivation(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return apperror.BadRequest("missing email")
	}

	if err := ctrl.svc.RequestActivation(email); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", email),
		nil,
	))
}

func (ctrl *controller) ChangePassword(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*ChangePasswordRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	userId := fmt.Sprintf("%v", c.Locals("userId"))

	if err := ctrl.svc.ChangePassword(userId, reqBody); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to change password",
		nil,
	))
}

func (ctrl *controller) RefreshAccessToken(c *fiber.Ctx) error {
	memberId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidUserSession)
	}

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidCompanySession)
	}

	roleId, err := helper.InterfaceToUint(c.Locals("tierLevel"))
	if err != nil {
		return apperror.Unauthorized("invalid tier level session")
	}

	apiKey := fmt.Sprintf("%v", c.Locals("apiKey"))

	accessToken, err := ctrl.svc.RefreshAccessToken(memberId, companyId, roleId, apiKey)
	if err != nil {
		return err
	}

	if err := setTokenCookie(c, "aif_token", accessToken, ctrl.cfg.Env.JwtExpiresMinutes); err != nil {
		return apperror.Internal("failed to set access token cookie", err)
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"access token refreshed",
		nil,
	))
}

func (ctrl *controller) Login(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*userLoginRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	accessToken, refreshToken, loginResp, err := ctrl.svc.LoginMember(reqBody)
	if err != nil {
		return err
	}

	const accessCookieName = "aif_token"
	const refreshCookieName = "aif_refresh_token"

	// Set access token cookie
	if err := setTokenCookie(c, accessCookieName, accessToken, ctrl.cfg.Env.JwtExpiresMinutes); err != nil {
		return apperror.Internal("failed to set access token cookie", err)
	}

	// Set refresh token cookie
	if err := setTokenCookie(c, refreshCookieName, refreshToken, ctrl.cfg.Env.JwtRefreshTokenExpiresMinutes); err != nil {
		return apperror.Internal("failed to set refresh token cookie", err)
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess("succeed to login", loginResp))
}

func (ctrl *controller) RequestPasswordReset(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*RequestPasswordResetRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	if err := ctrl.svc.RequestPasswordReset(reqBody.Email); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to reset your password", reqBody.Email),
		nil,
	))
}

func (ctrl *controller) PasswordReset(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*PasswordResetRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	token := c.Params("token")
	if token == "" {
		return apperror.BadRequest("missing password reset token")
	}

	if err := ctrl.svc.PasswordReset(token, reqBody); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to reset password",
		nil,
	))
}

func setTokenCookie(c *fiber.Ctx, name, value, durationStr string) error {
	minutes, err := strconv.Atoi(durationStr)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(minutes) * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return nil
}

func clearAuthCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax", // Atau "Strict" jika lebih aman
	})
}
