package middleware

import (
	"front-office/helper"
	"front-office/pkg/permission"
	"front-office/pkg/role"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

func IsRoleRequestValid(c *fiber.Ctx) error {
	request := &role.RoleRequest{}
	if err := c.BodyParser(&request); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		resp := helper.ResponseFailed(errValid.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	c.Locals("request", request)

	return c.Next()
}

func IsPermissionRequestValid(c *fiber.Ctx) error {
	request := &permission.PermissionRequest{}
	if err := c.BodyParser(&request); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		resp := helper.ResponseFailed(errValid.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	c.Locals("request", request)

	return c.Next()
}

func IsRegisterUserRequestValid(c *fiber.Ctx) error {
	request := &user.RegisterUserRequest{}
	if err := c.BodyParser(&request); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		resp := helper.ResponseFailed(errValid.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	c.Locals("request", request)

	return c.Next()
}

func IsLoginRequestValid(c *fiber.Ctx) error {
	request := &user.UserLoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		resp := helper.ResponseFailed(errValid.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	c.Locals("request", request)

	return c.Next()
}
