package middleware

import (
	"front-office/helper"
	"front-office/pkg/role"

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
