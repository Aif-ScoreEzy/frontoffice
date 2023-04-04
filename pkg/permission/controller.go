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
		ID:   result.ID,
		Name: result.Name,
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a permission by ID",
		dataRespose,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdatePermissionByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*PermissionRequest)
	id := c.Params("id")

	_, err := GetPermissionByIDSvc(id)
	if err != nil && err.Error() == "record not found" {
		resp := helper.ResponseFailed("Data is not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	_, err = GetPermissionByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	result, err := UpdatePermissionByIDSvc(*req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to update a permission",
		result,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeletePermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := GetPermissionByIDSvc(id)
	if err != nil && err.Error() == "record not found" {
		resp := helper.ResponseFailed("Data is not found")

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
