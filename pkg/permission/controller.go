package permission

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func CreatePermission(c *fiber.Ctx) error {
	req := c.Locals("request").(*PermissionRequest)

	_, err := GetPermissionByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	permission, err := CreatePermissionSvc(*req)
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

func GetPermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	permission, err := IsPermissionExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a permission by ID",
		permission,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdatePermissionByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*PermissionRequest)
	id := c.Params("id")

	_, err := IsPermissionExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	_, err = GetPermissionByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	permission, err := UpdatePermissionByIDSvc(*req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to update a permission",
		permission,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeletePermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := IsPermissionExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if err := DeletePermissionByIDSvc(id); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to delete a permission",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
