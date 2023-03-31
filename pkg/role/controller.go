package role

import (
	"github.com/gofiber/fiber/v2"
)

func CreateRole(c *fiber.Ctx) error {
	var roleRequest RoleRequest

	if err := c.BodyParser(&roleRequest); err != nil {
		return err
	}

	if errCreateRole := CreateRoleSvc(roleRequest); errCreateRole != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errCreateRole.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := GetRoleByIDSvc(id)
	if err != nil && err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data is not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	dataRespose := RoleResponse{
		Name: result.Name,
	}

	return c.Status(200).JSON(dataRespose)
}
