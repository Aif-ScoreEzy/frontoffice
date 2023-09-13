package user

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/company"
	"front-office/pkg/role"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterUserRequest)

	user, err := FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user != nil {
		resp := helper.ResponseFailed("Email already exists")
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, err = RegisterUserSvc(req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Active:  user.Active,
		Company: user.Company,
		Role:    user.Role,
	}

	resp := helper.ResponseSuccess(
		"Success to register",
		dataResponse,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func SendEmailVerification(c *fiber.Ctx) error {
	req := c.Locals("request").(*SendEmailVerificationRequest)

	user, err := FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user == nil {
		resp := helper.ResponseFailed("User not found")
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	err = SendEmailVerificationSvc(req, user)
	if err != nil {
		resp := helper.ResponseFailed("Something goes to wrong. Please try again")
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"The verification link has been sent to your email address",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func VerifyUser(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	user, err := FindOneByID(userID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user.IsVerified {
		statusCode, resp := helper.GetError(constant.AlreadyVerified)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = VerifyUserSvc(userID)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Your email has been verified",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func RequestPasswordReset(c *fiber.Ctx) error {
	req := c.Locals("request").(*RequestPasswordResetRequest)

	user, err := FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	err = SendEmailPasswordResetSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we sent email to %s with a link to reset your password", req.Email),
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func PasswordReset(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	req := c.Locals("request").(*PasswordResetRequest)

	_, err := PasswordResetSvc(userID, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Successful password reset",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func RegisterMember(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterMemberRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	loggedUser, err := FindUserByIDSvc(userID)
	// Only user with role of admin can register another user
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if loggedUser.Role.Name != "admin" {
		statusCode, resp := helper.GetError(constant.RequestProhibited)
		return c.Status(statusCode).JSON(resp)
	}

	user, err := FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = RegisterMemberSvc(req, loggedUser)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"We sent an email with the credentials for login to "+req.Email,
		nil,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func Login(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	user, err := FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user == nil {
		statusCode, resp := helper.GetError(constant.InvalidEmailOrPassword)
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
		RoleID:    user.RoleID,
		Token:     token,
	}

	resp := helper.ResponseSuccess(
		"Succeed to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func ChangePassword(c *fiber.Ctx) error {
	req := c.Locals("request").(*ChangePasswordRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	user, err := FindUserByIDSvc(userID)
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
		"Succeed to change password",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func ActivateUser(c *fiber.Ctx) error {
	key := c.Params("key")

	user, _ := FindUserByKeySvc(key)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)

		return c.Status(statusCode).JSON(resp)
	}

	user, err := ActivateUserByKeySvc(key)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success in activating the user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeactiveUser(c *fiber.Ctx) error {
	email := c.Params("email")

	user, _ := FindUserByEmailSvc(email)
	if user == nil {
		resp := helper.ResponseFailed("User not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	user, err := DeactivateUserByEmailSvc(email)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success in deactivating the user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateUserByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateUserRequest)
	id := c.Params("id")

	_, err := FindUserByIDSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, _ := FindUserByEmailSvc(req.Email)
	if user != nil {
		resp := helper.ResponseFailed("Email already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	_, err = company.IsCompanyIDExistSvc(req.CompanyID)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	_, err = role.IsRoleIDExistSvc(req.RoleID)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, err = UpdateUserByIDSvc(req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"Success to update user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetAllUsers(c *fiber.Ctx) error {
	users, err := GetAllUsersSvc()
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get all users",
		users,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, _ := FindUserByIDSvc(id)
	if user == nil {
		resp := helper.ResponseFailed("User not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a user by ID",
		user,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
