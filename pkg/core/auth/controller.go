package auth

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/passwordresettoken"
	"front-office/pkg/core/user"
	"front-office/utility/mailjet"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	svcUser user.Service,
	svcActivationToken activationtoken.Service,
	svcPasswordResetToken passwordresettoken.Service,
	cfg *config.Config,
) Controller {
	return &controller{
		Svc:                   service,
		SvcUser:               svcUser,
		SvcActivationToken:    svcActivationToken,
		SvcPasswordResetToken: svcPasswordResetToken,
		Cfg:                   cfg,
	}
}

type controller struct {
	Svc                   Service
	SvcUser               user.Service
	SvcActivationToken    activationtoken.Service
	SvcPasswordResetToken passwordresettoken.Service
	Cfg                   *config.Config
}

type Controller interface {
	RegisterAdmin(c *fiber.Ctx) error
	RegisterMemberAifCore(c *fiber.Ctx) error
	VerifyUser(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	SendEmailActivation(c *fiber.Ctx) error
	RequestPasswordResetAifCore(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	LoginAifCore(c *fiber.Ctx) error
	PasswordResetAifCore(c *fiber.Ctx) error
	// ChangePasswordAifcore(c *fiber.Ctx) error
}

func (ctrl *controller) RegisterAdmin(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterAdminRequest)

	userExists, _ := ctrl.SvcUser.FindUserByEmailSvc(req.Email)
	if userExists != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	newUser, token, err := ctrl.Svc.RegisterAdminSvc(req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	err = mailjet.SendEmailVerification(req.Email, token)
	if err != nil {
		resend := "resend"
		req := &user.UpdateUserRequest{
			Status: &resend,
		}

		_, err = ctrl.SvcUser.UpdateUserByIdSvc(req, newUser)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := RegisterAdminResponse{
		Id:      newUser.Id,
		Name:    newUser.Name,
		Email:   newUser.Email,
		Phone:   newUser.Phone,
		Status:  newUser.Status,
		Active:  newUser.Active,
		Company: newUser.Company,
		Role:    newUser.Role,
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", req.Email),
		dataResponse,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (ctrl *controller) RegisterMemberAifCore(c *fiber.Ctx) error {
	req := c.Locals("request").(*user.RegisterMemberRequest)

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resAddMember, err := ctrl.Svc.AddMemberAifCoreService(req, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	memberRoleId := 2
	token, err := ctrl.SvcActivationToken.CreateActivationTokenAifCore(resAddMember.Data.MemberId, companyId, uint(memberRoleId))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailActivation(req.Email, token)
	if err != nil {
		resend := "resend"
		req := &user.UpdateUserRequest{
			Status: &resend,
		}

		err = ctrl.SvcUser.UpdateUserByIdAifCore(req, resAddMember.Data.MemberId)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
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

	minutesToExpired, _ := strconv.Atoi(ctrl.Cfg.Env.JwtActivationExpiresMinutes)

	activationToken, err := ctrl.SvcActivationToken.FindActivationTokenByTokenSvc(token)
	if err != nil || (activationToken != nil && activationToken.Activation) {
		statusCode, resp := helper.GetError(constant.InvalidActivationLink)
		return c.Status(statusCode).JSON(resp)
	}

	userExists, err := ctrl.SvcUser.FindUserByIdSvc(activationToken.UserId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if userExists.IsVerified && userExists.Active {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	if activationToken != nil && time.Since(activationToken.CreatedAt).Minutes() > float64(minutesToExpired) {
		resend := "resend"
		req := &user.UpdateUserRequest{
			Status: &resend,
		}

		_, err = ctrl.SvcUser.UpdateUserByIdSvc(req, userExists)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.ActivationTokenExpired)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.Svc.VerifyUserTxSvc(userExists.Id, token, req)
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

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",              // Empty value
		Expires:  time.Unix(0, 0), // Expired time (epoch)
		HTTPOnly: true,            // HTTPOnly for security
		Secure:   true,
		SameSite: "Lax", // Adjust as needed
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",              // Empty value
		Expires:  time.Unix(0, 0), // Expired time (epoch)
		HTTPOnly: true,            // HTTPOnly for security
		Secure:   true,
		SameSite: "Lax", // Adjust as needed
	})

	resp := helper.ResponseSuccess(
		"succeed to logout",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) SendEmailActivation(c *fiber.Ctx) error {
	email := c.Params("email")

	userExists, err := ctrl.SvcUser.FindUserAifCore(&user.FindUserQuery{
		Email: email,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if userExists.Data.IsVerified {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	token, err := ctrl.SvcActivationToken.CreateActivationTokenAifCore(userExists.Data.MemberId, userExists.Data.CompanyId, userExists.Data.RoleId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailActivation(email, token)
	if err != nil {
		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	} else {
		pending := "pending"
		req := &user.UpdateUserRequest{
			Status: &pending,
		}

		err = ctrl.SvcUser.UpdateUserByIdAifCore(req, userExists.Data.MemberId)
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
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	userExists, err := ctrl.SvcUser.FindUserByIdAndCompanyIdSvc(userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.Svc.ChangePasswordSvc(userExists, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
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

	secret := ctrl.Cfg.Env.JwtSecretKey
	accessTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtExpiresMinutes)
	newAccessToken, err := helper.GenerateToken(secret, accessTokenExpirationMinutes, userId, companyId, uint(tierLevel))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
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

func (ctrl *controller) LoginAifCore(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	res, err := ctrl.SvcUser.FindUserAifCore(&user.FindUserQuery{
		Email: req.Email,
	})
	if err != nil {
		resp := helper.ResponseFailed(err.Error())
		return c.Status(res.StatusCode).JSON(resp)
	}

	if res == nil {
		statusCode, resp := helper.GetError(constant.InvalidEmailOrPassword)
		return c.Status(statusCode).JSON(resp)
	} else if res != nil && !res.Data.Active {
		statusCode, resp := helper.GetError(constant.RequestProhibited)
		return c.Status(statusCode).JSON(resp)
	}

	accessToken, refreshToken, err := ctrl.Svc.LoginMemberAifCoreService(req, res.Data)
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
		Name:     "aif_refreh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(refreshTokenExpirationMinutes) * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	if res.StatusCode != 200 {
		resp := helper.ResponseFailed(res.Message)
		return c.Status(res.StatusCode).JSON(resp)
	}

	user := res.Data
	data := UserLoginResponse{
		Id:          user.MemberId,
		Name:        user.Name,
		Email:       user.Email,
		CompanyId:   user.CompanyId,
		CompanyName: user.MstCompany.CompanyName,
		TierLevel:   user.RoleId,
		Image:       user.Image,
	}

	responseSuccess := helper.ResponseSuccess(
		"succeed to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(responseSuccess)
}

func (ctrl *controller) RequestPasswordResetAifCore(c *fiber.Ctx) error {
	req := c.Locals("request").(*RequestPasswordResetRequest)

	userExists, err := ctrl.SvcUser.FindUserAifCore(&user.FindUserQuery{
		Email: req.Email,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if !userExists.Data.IsVerified {
		statusCode, resp := helper.GetError(constant.UnverifiedUser)
		return c.Status(statusCode).JSON(resp)
	}

	token, err := ctrl.SvcPasswordResetToken.CreatePasswordResetTokenAifCore(userExists.Data.MemberId, userExists.Data.CompanyId, userExists.Data.RoleId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailPasswordReset(req.Email, userExists.Data.Name, token)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to reset your password", req.Email),
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) PasswordResetAifCore(c *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	result, err := ctrl.SvcPasswordResetToken.FindPasswordResetTokenByTokenSvc(token)
	if err != nil || result == nil || result.Data == nil || result.Data.Activation {
		statusCode, resp := helper.GetError(constant.InvalidPasswordResetLink)
		return c.Status(statusCode).JSON(resp)
	}

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

	err = ctrl.Svc.PasswordResetSvc(userId, token, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to reset password",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
