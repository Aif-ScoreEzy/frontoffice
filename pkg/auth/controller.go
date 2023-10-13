package auth

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/user"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdmin(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterAdminRequest)

	user, _ := user.FindUserByEmailSvc(req.Email)
	if user != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	user, err := RegisterAdminSvc(req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	sendVerificationRequest := &SendEmailVerificationRequest{
		Email: req.Email,
	}

	err = SendEmailVerificationSvc(sendVerificationRequest, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := RegisterAdminResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Active:  user.Active,
		Company: user.Company,
		Role:    user.Role,
	}

	resp := helper.ResponseSuccess(
		"the verification link has been sent to your email address",
		dataResponse,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func SendEmailVerification(c *fiber.Ctx) error {
	req := c.Locals("request").(*SendEmailVerificationRequest)

	user, err := user.FindUserByEmailSvc(req.Email)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	} else if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = SendEmailVerificationSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"the verification link has been sent to your email address",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func VerifyUser(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	data, _ := VerifyActivationToken(token)
	if data.Activation || data == nil {
		statusCode, resp := helper.GetError(constant.InvalidActivationLink)
		return c.Status(statusCode).JSON(resp)
	}

	result, err := user.FindOneByID(userID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if result.IsVerified && result.Active {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = VerifyUserTxSvc(userID, token, req)
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

func RequestPasswordReset(c *fiber.Ctx) error {
	req := c.Locals("request").(*RequestPasswordResetRequest)

	user, err := user.FindUserByEmailSvc(req.Email)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	} else if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = SendEmailPasswordResetSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("We've sent an email to %s with a link to reset your password", req.Email),
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func PasswordReset(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	req := c.Locals("request").(*PasswordResetRequest)
	token := c.Params("token")

	data, _ := VerifyPasswordResetToken(token)
	if data.Activation || data == nil {
		statusCode, resp := helper.GetError(constant.InvalidPasswordResetLink)
		return c.Status(statusCode).JSON(resp)
	}

	err := PasswordResetSvc(userID, token, req)
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

func Login(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	user, err := user.FindUserByEmailSvc(req.Email)
	if user == nil {
		statusCode, resp := helper.GetError(constant.InvalidEmailOrPassword)
		return c.Status(statusCode).JSON(resp)
	} else if !user.Active {
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
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CompanyID: user.CompanyID,
		TierLevel: user.Role.TierLevel,
		Token:     token,
	}

	resp := helper.ResponseSuccess(
		"succeed to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func ChangePassword(c *fiber.Ctx) error {
	req := c.Locals("request").(*ChangePasswordRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	user, err := user.FindUserByIDSvc(userID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ChangePasswordSvc(userID, user, req)
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

func UpdateProfile(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateProfileRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	updateUser, err := UpdateProfileSvc(userID, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := user.UserUpdateResponse{
		ID:        updateUser.ID,
		Name:      updateUser.Name,
		Email:     updateUser.Email,
		Active:    updateUser.Active,
		CompanyID: updateUser.CompanyID,
		RoleID:    updateUser.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success to update user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UploadProfileImage(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	file, err := c.FormFile("image")
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if file.Size > 200*1024 { // 200 KB dalam byte
		statusCode, resp := helper.GetError(constant.FileSizeIsTooLarge)
		return c.Status(statusCode).JSON(resp)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", userID, ext)
	filePath := fmt.Sprintf("./public/%s", filename)

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	if err := c.SaveFile(file, filePath); err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	updateUser, err := UploadProfileImageSvc(userID, &filename)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := user.UserUpdateResponse{
		ID:        updateUser.ID,
		Name:      updateUser.Name,
		Email:     updateUser.Email,
		Active:    updateUser.Active,
		CompanyID: updateUser.CompanyID,
		RoleID:    updateUser.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success to upload profile image",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
