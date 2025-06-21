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
	"front-office/utility/mailjet"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
		Svc:                   service,
		SvcUser:               svcUser,
		SvcActivationToken:    svcActivationToken,
		SvcPasswordResetToken: svcPasswordResetToken,
		SvcLogOperation:       svcLogOperation,
		Cfg:                   cfg,
	}
}

type controller struct {
	Svc                   Service
	SvcUser               member.Service
	SvcActivationToken    activationtoken.Service
	SvcPasswordResetToken passwordresettoken.Service
	SvcLogOperation       operation.Service
	Cfg                   *config.Config
}

type Controller interface {
	RegisterMember(c *fiber.Ctx) error
	VerifyUser(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	SendEmailActivation(c *fiber.Ctx) error
	RequestPasswordReset(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	PasswordReset(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
}

func (ctrl *controller) RegisterMember(c *fiber.Ctx) error {
	req, ok := c.Locals("request").(*member.RegisterMemberRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	currentUserId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		return apperror.Unauthorized("invalid user session")
	}

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		return apperror.Unauthorized("invalid company session")
	}

	req.CompanyId = companyId
	req.RoleId = uint(memberRoleId)

	if err := ctrl.Svc.AddMember(currentUserId, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", req.Email),
		nil,
	))
}

func (ctrl *controller) VerifyUser(c *fiber.Ctx) error {
	req, ok := c.Locals("request").(*PasswordResetRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	token := c.Params("token")
	if token == "" {
		return apperror.BadRequest("missing activation token")
	}

	if err := ctrl.Svc.VerifyMember(token, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"your account has been verified",
		nil,
	))
}

func (ctrl *controller) Logout(c *fiber.Ctx) error {
	memberId, _ := helper.InterfaceToUint(c.Locals("userId"))
	companyId, _ := helper.InterfaceToUint(c.Locals("companyId"))

	c.Cookie(&fiber.Cookie{
		Name:     "aif_token",
		Value:    "",              // Empty value
		Expires:  time.Unix(0, 0), // Expired time (epoch)
		HTTPOnly: true,            // HTTPOnly for security
		Secure:   true,
		SameSite: "Lax", // Adjust as needed
	})

	c.Cookie(&fiber.Cookie{
		Name:     "aif_refresh_token",
		Value:    "",              // Empty value
		Expires:  time.Unix(0, 0), // Expired time (epoch)
		HTTPOnly: true,            // HTTPOnly for security
		Secure:   true,
		SameSite: "Lax", // Adjust as needed
	})

	addLogRequest := &operation.AddLogRequest{
		MemberId:  memberId,
		CompanyId: companyId,
		Action:    constant.EventSignOut,
	}

	err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for user logout")
	}

	resp := helper.ResponseSuccess(
		"succeed to logout",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) SendEmailActivation(c *fiber.Ctx) error {
	email := c.Params("email")

	memberData, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
		Email: email,
	})
	if err != nil {
		return err
	}

	if memberData.MemberId == 0 {
		return err
	}

	if memberData.IsVerified {
		return err
	}

	token, err := ctrl.SvcActivationToken.CreateActivationToken(memberData.MemberId, memberData.CompanyId, memberData.RoleId)
	if err != nil {
		return err
	}

	memberId := fmt.Sprintf("%d", memberData.MemberId)
	err = mailjet.SendEmailActivation(email, token)
	if err != nil {
		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	} else {
		pending := "pending"
		req := &member.UpdateUserRequest{
			MailStatus: &pending,
		}

		err = ctrl.SvcUser.UpdateMemberById(memberId, req)
		if err != nil {
			return err
		}
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", email),
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) ChangePassword(c *fiber.Ctx) error {
	req := c.Locals("request").(*ChangePasswordRequest)
	memberId := fmt.Sprintf("%v", c.Locals("userId"))

	memberData, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
		Id: memberId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	result, err := ctrl.Svc.ChangePassword(memberId, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if !result.Success {
		statusCode, resp := helper.GetError(result.Message)
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendConfirmationEmailPasswordChangeSuccess(memberData.Name, memberData.Email)
	if err != nil {
		return err
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  memberData.MemberId,
		CompanyId: memberData.CompanyId,
		Action:    constant.EventChangePassword,
	}

	err = ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for change password")
	}

	resp := helper.ResponseSuccess(
		"succeed to change password",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) RefreshAccessToken(c *fiber.Ctx) error {
	userId, _ := helper.InterfaceToUint(c.Locals("userId"))
	companyId, _ := helper.InterfaceToUint(c.Locals("companyId"))
	tierLevel, _ := helper.InterfaceToUint(c.Locals("tierLevel"))
	apiKey, _ := c.Locals("apiKey").(string)

	secret := ctrl.Cfg.Env.JwtSecretKey
	accessTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtExpiresMinutes)
	newAccessToken, err := helper.GenerateToken(secret, accessTokenExpirationMinutes, userId, companyId, uint(tierLevel), apiKey)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "aif_token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(time.Duration(accessTokenExpirationMinutes) * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	resp := helper.ResponseSuccess(
		"access token refreshed",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) Login(c *fiber.Ctx) error {
	reqBody, ok := c.Locals("request").(*userLoginRequest)
	if !ok {
		return apperror.BadRequest(constant.InvalidRequestFormat)
	}

	accessToken, refreshToken, loginResp, err := ctrl.Svc.LoginMember(reqBody)
	if err != nil {
		return err
	}

	const accessCookieName = "aif_token"
	const refreshCookieName = "aif_refresh_token"

	// Set access token cookie
	if err := ctrl.setTokenCookie(c, accessCookieName, accessToken, ctrl.Cfg.Env.JwtExpiresMinutes); err != nil {
		return apperror.Internal("failed to set access token cookie", err)
	}

	// Set refresh token cookie
	if err := ctrl.setTokenCookie(c, refreshCookieName, refreshToken, ctrl.Cfg.Env.JwtRefreshTokenExpiresMinutes); err != nil {
		return apperror.Internal("failed to set refresh token cookie", err)
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess("succeed to login", loginResp))
}

func (ctrl *controller) RequestPasswordReset(c *fiber.Ctx) error {
	req := c.Locals("request").(*RequestPasswordResetRequest)

	memberData, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
		Email: req.Email,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if memberData.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.UserNotFoundForgotEmail)
		return c.Status(statusCode).JSON(resp)
	}

	if !memberData.IsVerified {
		statusCode, resp := helper.GetError(constant.UnverifiedUser)
		return c.Status(statusCode).JSON(resp)
	}

	token, err := ctrl.SvcPasswordResetToken.CreatePasswordResetTokenAifCore(memberData.MemberId, memberData.CompanyId, memberData.RoleId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailPasswordReset(req.Email, memberData.Name, token)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  memberData.MemberId,
		CompanyId: memberData.CompanyId,
		Action:    constant.EventRequestPasswordReset,
	}

	err = ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for request password reset")
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to reset your password", req.Email),
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) PasswordReset(c *fiber.Ctx) error {
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	result, err := ctrl.SvcPasswordResetToken.FindPasswordResetTokenByTokenSvc(token)
	if err != nil || result == nil || result.Data == nil {
		statusCode, resp := helper.GetError(constant.InvalidPasswordResetLink)
		return c.Status(statusCode).JSON(resp)
	}

	memberId := result.Data.Member.MemberId
	companyId := result.Data.Member.CompanyId

	jwtResetPasswordExpiresMinutesStr := ctrl.Cfg.Env.JwtResetPasswordExpiresMinutes
	minutesToExpired, err := strconv.Atoi(jwtResetPasswordExpiresMinutesStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	elapsedMinutes := time.Since(result.Data.CreatedAt).Minutes()
	if elapsedMinutes > float64(minutesToExpired) {
		_, err := ctrl.SvcPasswordResetToken.DeletePasswordResetToken(result.Data.Id)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.InvalidPasswordResetLink)
		return c.Status(statusCode).JSON(resp)
	}

	err = ctrl.Svc.PasswordResetSvc(memberId, token, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  memberId,
		CompanyId: companyId,
		Action:    constant.EventPasswordReset,
	}

	err = ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for password reset")
	}

	resp := helper.ResponseSuccess(
		"succeed to reset password",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) setTokenCookie(c *fiber.Ctx, name, value, durationStr string) error {
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
