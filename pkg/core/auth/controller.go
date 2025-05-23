package auth

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
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
	req := c.Locals("request").(*member.RegisterMemberRequest)

	currentUserId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	memberRoleId := 2
	req.CompanyId = companyId
	req.RoleId = uint(memberRoleId)
	resAddMember, err := ctrl.Svc.AddMember(req, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if resAddMember != nil && resAddMember.StatusCode != 200 {
		_, resp := helper.GetError(resAddMember.Message)
		return c.Status(resAddMember.StatusCode).JSON(resp)
	}

	token, result, err := ctrl.SvcActivationToken.CreateActivationToken(resAddMember.Data.MemberId, companyId, uint(memberRoleId))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if !result.Success {
		statusCode, resp := helper.GetError(fiber.ErrInternalServerError.Error())

		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailActivation(req.Email, token)
	if err != nil {
		resend := "resend"
		req := &member.UpdateUserRequest{
			MailStatus: &resend,
		}

		memberId := fmt.Sprintf("%d", result.Data.MemberId)
		_, err := ctrl.SvcUser.UpdateMemberById(memberId, req)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  currentUserId,
		CompanyId: companyId,
		Action:    constant.EventRegisterMember,
	}

	resAddLog, err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil || !resAddLog.Success {
		log.Println("Failed to log operation for register member")
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", req.Email),
		nil,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (ctrl *controller) VerifyUser(c *fiber.Ctx) error {
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	result, err := ctrl.SvcActivationToken.FindActivationTokenByToken(token)
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(constant.InvalidActivationLink)
		return c.Status(statusCode).JSON(resp)
	}

	memberId := fmt.Sprintf("%d", result.Data.MemberId)

	userExists, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
		Id: memberId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if userExists.Data.IsVerified && userExists.Data.Active {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	minutesToExpired, err := strconv.Atoi(ctrl.Cfg.Env.JwtActivationExpiresMinutes)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	elapsedMinutes := time.Since(result.Data.CreatedAt).Minutes()
	if elapsedMinutes > float64(minutesToExpired) {
		resend := "resend"
		req := &member.UpdateUserRequest{
			MailStatus: &resend,
		}

		_, err = ctrl.SvcUser.UpdateMemberById(memberId, req)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.ActivationTokenExpired)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.Svc.VerifyMemberAif(userExists.Data.MemberId, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"your account has been verified",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
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

	resAddLog, err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil || !resAddLog.Success {
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

	userExists, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
		Email: email,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if userExists.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	if userExists.Data.IsVerified {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	token, result, err := ctrl.SvcActivationToken.CreateActivationToken(userExists.Data.MemberId, userExists.Data.CompanyId, userExists.Data.RoleId)
	if err != nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	memberId := fmt.Sprintf("%d", userExists.Data.MemberId)
	err = mailjet.SendEmailActivation(email, token)
	if err != nil {
		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	} else {
		pending := "pending"
		req := &member.UpdateUserRequest{
			MailStatus: &pending,
		}

		_, err = ctrl.SvcUser.UpdateMemberById(memberId, req)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
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

	member, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
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

	err = mailjet.SendConfirmationEmailPasswordChangeSuccess(member.Data.Name, member.Data.Email)
	if err != nil {
		return err
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  member.Data.MemberId,
		CompanyId: member.Data.CompanyId,
		Action:    constant.EventChangePassword,
	}

	resAddLog, err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil || !resAddLog.Success {
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
	req := c.Locals("request").(*UserLoginRequest)

	res, err := ctrl.Svc.LoginMember(req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if res != nil && !res.Success {
		_, resp := helper.GetError(res.Message)
		return c.Status(res.StatusCode).JSON(resp)
	}

	accessToken, refreshToken, err := ctrl.Svc.generateTokens(res.Data.MemberId, res.Data.CompanyId, res.Data.RoleId, res.Data.ApiKey)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	accessTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtExpiresMinutes)
	c.Cookie(&fiber.Cookie{
		Name:     "aif_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Duration(accessTokenExpirationMinutes) * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	refreshTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtRefreshTokenExpiresMinutes)
	c.Cookie(&fiber.Cookie{
		Name:     "aif_refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(refreshTokenExpirationMinutes) * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	addLogRequest := &operation.AddLogRequest{
		MemberId:  res.Data.MemberId,
		CompanyId: res.Data.CompanyId,
		Action:    constant.EventSignIn,
	}

	resAddLog, err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil || !resAddLog.Success {
		log.Println("Failed to log operation for user login")
	}

	data := &UserLoginResponse{
		Id:                 res.Data.MemberId,
		Name:               res.Data.Name,
		Email:              res.Data.Email,
		CompanyId:          res.Data.CompanyId,
		CompanyName:        res.Data.CompanyName,
		TierLevel:          res.Data.RoleId,
		Image:              res.Data.Image,
		SubscriberProducts: res.Data.SubscriberProducts,
	}

	responseSuccess := helper.ResponseSuccess(
		"succeed to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(responseSuccess)
}

func (ctrl *controller) RequestPasswordReset(c *fiber.Ctx) error {
	req := c.Locals("request").(*RequestPasswordResetRequest)

	member, err := ctrl.SvcUser.GetMemberBy(&member.FindUserQuery{
		Email: req.Email,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if member.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.UserNotFoundForgotEmail)
		return c.Status(statusCode).JSON(resp)
	}

	if !member.Data.IsVerified {
		statusCode, resp := helper.GetError(constant.UnverifiedUser)
		return c.Status(statusCode).JSON(resp)
	}

	token, err := ctrl.SvcPasswordResetToken.CreatePasswordResetTokenAifCore(member.Data.MemberId, member.Data.CompanyId, member.Data.RoleId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailPasswordReset(req.Email, member.Data.Name, token)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  member.Data.MemberId,
		CompanyId: member.Data.CompanyId,
		Action:    constant.EventRequestPasswordReset,
	}

	resAddLog, err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil || !resAddLog.Success {
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

	resAddLog, err := ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil || !resAddLog.Success {
		log.Println("Failed to log operation for password reset")
	}

	resp := helper.ResponseSuccess(
		"succeed to reset password",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
