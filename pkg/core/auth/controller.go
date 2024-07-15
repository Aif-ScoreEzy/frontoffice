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
	RegisterMember(c *fiber.Ctx) error
	VerifyUser(c *fiber.Ctx) error
	// Login(c *fiber.Ctx) error
	// SendEmailActivation(c *fiber.Ctx) error
	// RequestPasswordReset(c *fiber.Ctx) error
	PasswordReset(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	// RefreshAccessToken(c *fiber.Ctx) error
	LoginAifCore(c *fiber.Ctx) error
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

		_, err = ctrl.SvcUser.UpdateUserByIDSvc(req, newUser)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := RegisterAdminResponse{
		ID:      newUser.ID,
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

func (ctrl *controller) RegisterMember(c *fiber.Ctx) error {
	req := c.Locals("request").(*user.RegisterMemberRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	userExists, _ := ctrl.SvcUser.FindUserByEmailSvc(req.Email)
	if userExists != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	result, token, err := ctrl.Svc.RegisterMemberSvc(req, companyID)
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

		_, err = ctrl.SvcUser.UpdateUserByIDSvc(req, result)
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

	userExists, err := ctrl.SvcUser.FindUserByIDSvc(activationToken.UserID)
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

		_, err = ctrl.SvcUser.UpdateUserByIDSvc(req, userExists)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.ActivationTokenExpired)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.Svc.VerifyUserTxSvc(userExists.ID, token, req)
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

// func (ctrl *controller) Login(c *fiber.Ctx) error {
// 	req := c.Locals("request").(*UserLoginRequest)

// 	user, err := ctrl.SvcUser.FindUserByEmailSvc(req.Email)
// 	if user == nil {
// 		statusCode, resp := helper.GetError(constant.InvalidEmailOrPassword)
// 		return c.Status(statusCode).JSON(resp)
// 	} else if user != nil && !user.Active {
// 		statusCode, resp := helper.GetError(constant.RequestProhibited)
// 		return c.Status(statusCode).JSON(resp)
// 	} else if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	accessToken, refreshToken, err := ctrl.Svc.LoginSvc(req, user)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	accessTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtExpiresMinutes)
// 	c.Cookie(&fiber.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		Expires:  time.Now().Add(time.Duration(accessTokenExpirationMinutes) * time.Minute),
// 		HTTPOnly: true,
// 		Secure:   true,
// 		SameSite: "Lax",
// 	})

// 	refreshTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtRefreshTokenExpiresMinutes)
// 	c.Cookie(&fiber.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		Expires:  time.Now().Add(time.Duration(refreshTokenExpirationMinutes) * time.Minute),
// 		HTTPOnly: true,
// 		Secure:   true,
// 		SameSite: "Lax",
// 	})

// 	data := UserLoginResponse{
// 		ID:          user.ID,
// 		Name:        user.Name,
// 		Email:       user.Email,
// 		CompanyID:   user.CompanyID,
// 		CompanyName: user.Company.CompanyName,
// 		TierLevel:   user.Role.TierLevel,
// 		Image:       user.Image,
// 	}

// 	resp := helper.ResponseSuccess(
// 		"succeed to login",
// 		data,
// 	)

// 	return c.Status(fiber.StatusOK).JSON(resp)
// }

// func (ctrl *controller) SendEmailActivation(c *fiber.Ctx) error {
// 	email := c.Params("email")
// 	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

// 	userExists, err := ctrl.SvcUser.FindUserByEmailSvc(email)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	userExists, err = ctrl.SvcUser.FindUserByIDAndCompanyIDSvc(userExists.ID, companyID)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	if userExists.IsVerified {
// 		statusCode, resp := helper.GetError(constant.AlreadyVerified)
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	token, _, err := ctrl.SvcActivationToken.CreateActivationTokenSvc(userExists.ID, userExists.CompanyID, userExists.Role.TierLevel)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	err = mailjet.SendEmailActivation(email, token)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(constant.SendEmailFailed)
// 		return c.Status(statusCode).JSON(resp)
// 	} else {
// 		pending := "pending"
// 		req := &user.UpdateUserRequest{
// 			Status: &pending,
// 		}

// 		_, err = ctrl.SvcUser.UpdateUserByIDSvc(req, userExists)
// 		if err != nil {
// 			statusCode, resp := helper.GetError(err.Error())
// 			return c.Status(statusCode).JSON(resp)
// 		}
// 	}

// 	resp := helper.ResponseSuccess(
// 		fmt.Sprintf("we've sent an email to %s with a link to activate the account", email),
// 		nil,
// 	)

// 	return c.Status(fiber.StatusOK).JSON(resp)
// }

// func (ctrl *controller) RequestPasswordReset(c *fiber.Ctx) error {
// 	req := c.Locals("request").(*RequestPasswordResetRequest)

// 	userExists, err := ctrl.SvcUser.FindUserByEmailSvc(req.Email)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	if !userExists.IsVerified {
// 		statusCode, resp := helper.GetError(constant.UnverifiedUser)
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	token, _, err := ctrl.SvcPasswordResetToken.CreatePasswordResetTokenSvc(userExists)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	err = mailjet.SendEmailPasswordReset(req.Email, userExists.Name, token)
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	resp := helper.ResponseSuccess(
// 		fmt.Sprintf("we've sent an email to %s with a link to reset your password", req.Email),
// 		nil,
// 	)

// 	return c.Status(fiber.StatusOK).JSON(resp)
// }

func (ctrl *controller) PasswordReset(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	data, err := ctrl.SvcPasswordResetToken.FindPasswordResetTokenByTokenSvc(token)
	if err != nil || (data != nil && data.Activation) {
		statusCode, resp := helper.GetError(constant.InvalidPasswordResetLink)
		return c.Status(statusCode).JSON(resp)
	}

	err = ctrl.Svc.PasswordResetSvc(userID, token, req)
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

func (ctrl *controller) ChangePassword(c *fiber.Ctx) error {
	req := c.Locals("request").(*ChangePasswordRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	userExists, err := ctrl.SvcUser.FindUserByIDAndCompanyIDSvc(userID, companyID)
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

// func (ctrl *controller) RefreshAccessToken(c *fiber.Ctx) error {
// 	userID := fmt.Sprintf("%v", c.Locals("userID"))
// 	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
// 	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)

// 	secret := ctrl.Cfg.Env.JwtSecretKey
// 	accessTokenExpirationMinutes, _ := strconv.Atoi(ctrl.Cfg.Env.JwtExpiresMinutes)
// 	newAccessToken, err := helper.GenerateToken(secret, accessTokenExpirationMinutes, userID, companyID, uint(tierLevel))
// 	if err != nil {
// 		statusCode, resp := helper.GetError(err.Error())
// 		return c.Status(statusCode).JSON(resp)
// 	}

// 	c.Cookie(&fiber.Cookie{
// 		Name:     "access_token",
// 		Value:    newAccessToken,
// 		Expires:  time.Now().Add(time.Duration(accessTokenExpirationMinutes) * time.Minute),
// 		HTTPOnly: true,
// 		Secure:   true,
// 		SameSite: "Lax",
// 	})

// 	resp := helper.ResponseSuccess(
// 		"access token refreshed",
// 		nil,
// 	)

// 	return c.Status(fiber.StatusOK).JSON(resp)
// }

func (ctrl *controller) LoginAifCore(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	res, err := ctrl.SvcUser.FindUserByEmailAifCore(req.Email)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())
		return c.Status(res.StatusCode).JSON(resp)
	}

	accessToken, refreshToken, err := ctrl.Svc.LoginAifCoreService(req, res.Data)
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
		ID:          user.MemberID,
		Name:        user.Name,
		Email:       user.Email,
		CompanyID:   user.CompanyId,
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
