package role

import (
	"front-office/helper"
	"front-office/pkg/permission"

	"github.com/gofiber/fiber/v2"
)

func CreateRole(c *fiber.Ctx) error {
	request := c.Locals("request").(*CreateRoleRequest)

	_, err := GetRoleByNameSvc(request.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	for _, permissionData := range request.Permissions {
		_, err := permission.IsPermissionExistSvc(permissionData.ID)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}
	}

	role, err := CreateRoleSvc(request)
	if err != nil && err.Error() == "record not found" {
		resp := helper.ResponseFailed("Data is not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	} else if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to create a role",
		role,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetAllRoles(c *fiber.Ctx) error {
	roles, err := GetAllRolesSvc()
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get all roles",
		roles,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	role, err := IsRoleIDExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a role by ID",
		role,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateRole(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateRoleRequest)
	id := c.Params("id")

	_, err := IsRoleIDExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	_, err = GetRoleByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	role, err := UpdateRoleByIDSvc(req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to update a role",
		role,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := IsRoleIDExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if err := DeleteRoleByIDSvc(id); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to delete a role",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
