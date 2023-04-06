package user

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterUserRequest)

	isUsernameExist, _ := IsUsernameExistSvc(req.Username)
	if isUsernameExist {
		resp := helper.ResponseFailed("Username already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	isEmailExist := IsEmailExistSvc(req.Email)
	if isEmailExist {
		resp := helper.ResponseFailed("Email already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, err := RegisterUserSvc(*req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to register",
		user,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func Login(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	isUsernameExist, user := IsUsernameExistSvc(req.Username)
	if !isUsernameExist {
		resp := helper.ResponseFailed("Username or password is incorrect")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	token, err := LoginSvc(*req, user)
	if err != nil && err.Error() == "password is incorrect" {
		resp := helper.ResponseFailed("Username or password is incorrect")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	} else if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	data := UserLoginResponse{
		Username:  user.Username,
		Email:     user.Email,
		CompanyID: user.CompanyID,
		Role:      user.RoleID,
		Key:       user.Key,
		Token:     token,
	}

	resp := helper.ResponseSuccess(
		"Success to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
