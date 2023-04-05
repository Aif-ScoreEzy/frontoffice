package user

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterUserRequest)

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
