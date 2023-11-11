package auth

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	activation_token "front-office/pkg/activation-token"
	"front-office/pkg/password_reset_token"
	"front-office/pkg/user"
	"front-office/utility/mailjet"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdmin(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterAdminRequest)

	userExists, _ := user.FindUserByEmailSvc(req.Email)
	if userExists != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	newUser, token, err := RegisterAdminSvc(req)
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

		_, err = user.UpdateUserByIDSvc(req, newUser)
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

func RegisterMember(c *fiber.Ctx) error {
	req := c.Locals("request").(*user.RegisterMemberRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	userExists, _ := user.FindUserByEmailSvc(req.Email)
	if userExists != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	result, token, err := RegisterMemberSvc(req, companyID)
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

		_, err = user.UpdateUserByIDSvc(req, result)
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

func VerifyUser(c *fiber.Ctx) error {
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_ACTIVATION_EXPIRES_MINUTES"))

	activationToken, err := activation_token.FindActivationTokenByTokenSvc(token)
	if err != nil || (activationToken != nil && activationToken.Activation) {
		statusCode, resp := helper.GetError(constant.InvalidActivationLink)
		return c.Status(statusCode).JSON(resp)
	}

	userExists, err := user.FindOneByUserID(activationToken.UserID)
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

		_, err = user.UpdateUserByIDSvc(req, userExists)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.InvalidActivationLink)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = VerifyUserTxSvc(userExists.ID, token, req)
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

func Login(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	user, err := user.FindUserByEmailSvc(req.Email)
	if user == nil {
		statusCode, resp := helper.GetError(constant.InvalidEmailOrPassword)
		return c.Status(statusCode).JSON(resp)
	} else if user != nil && !user.Active {
		statusCode, resp := helper.GetError(constant.RequestProhibited)
		return c.Status(statusCode).JSON(resp)
	} else if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	token, err := LoginSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	data := UserLoginResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CompanyID:   user.CompanyID,
		CompanyName: user.Company.CompanyName,
		TierLevel:   user.Role.TierLevel,
		Image:       user.Image,
		Token:       token,
	}

	resp := helper.ResponseSuccess(
		"succeed to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func SendEmailActivation(c *fiber.Ctx) error {
	email := c.Params("email")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	userExists, err := user.FindUserByEmailSvc(email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	userExists, err = user.FindUserByIDSvc(userExists.ID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if userExists.IsVerified {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	token, _, err := activation_token.CreateActivationTokenSvc(userExists.ID, userExists.CompanyID, userExists.Role.TierLevel)
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

		_, err = user.UpdateUserByIDSvc(req, userExists)
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

func RequestPasswordReset(c *fiber.Ctx) error {
	req := c.Locals("request").(*RequestPasswordResetRequest)

	userExists, err := user.FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	token, _, err := password_reset_token.CreatePasswordResetTokenSvc(userExists)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = mailjet.SendEmailPasswordReset(req.Email, userExists.Name, token)
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

func PasswordReset(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	data, err := password_reset_token.FindPasswordResetTokenByTokenSvc(token)
	if err != nil || (data != nil && data.Activation) {
		statusCode, resp := helper.GetError(constant.InvalidPasswordResetLink)
		return c.Status(statusCode).JSON(resp)
	}

	err = PasswordResetSvc(userID, token, req)
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

func ChangePassword(c *fiber.Ctx) error {
	req := c.Locals("request").(*ChangePasswordRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	userExists, err := user.FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ChangePasswordSvc(userExists, req)
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
