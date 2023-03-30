package role

import (
	"github.com/gofiber/fiber/v2"
)

func CreateRole(c *fiber.Ctx) error {
	var roleRequest RoleRequest

	if err := c.BodyParser(&roleRequest); err != nil {
		return err
	}

	if errCreateRole := CreateRoleService(roleRequest); errCreateRole != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errCreateRole.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
