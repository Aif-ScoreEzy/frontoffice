package permission

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func CreatePermission(c *fiber.Ctx) error {
	request := c.Locals("request").(*PermissionRequest)

	permission, err := CreatePermissionSvc(*request)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to create a permission",
		permission,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
