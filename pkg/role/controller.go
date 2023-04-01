package role

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func CreateRole(c *fiber.Ctx) error {
	request := c.Locals("request").(*RoleRequest)

	role, err := CreateRoleSvc(*request)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to create a role",
		role,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := GetRoleByIDSvc(id)
	if err != nil && err.Error() == "record not found" {
		resp := helper.ResponseFailed("Data is not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	} else if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataRespose := RoleResponse{
		Name: result.Name,
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a role by ID",
		dataRespose,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateRole(c *fiber.Ctx) error {
	request := c.Locals("request").(*RoleRequest)
	id := c.Params("id")

	_, err := GetRoleByIDSvc(id)
	if err != nil && err.Error() == "record not found" {
		resp := helper.ResponseFailed("Data is not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	result, err := UpdateRoleByIDSvc(*request, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to update a role",
		result,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
