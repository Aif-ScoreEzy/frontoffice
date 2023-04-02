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

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := GetPermissionByIDSvc(id)
	if err != nil && err.Error() == "record not found" {
		resp := helper.ResponseFailed("Data is not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	} else if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataRespose := PermissionResponse{
		Name: result.Name,
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a permission by ID",
		dataRespose,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
