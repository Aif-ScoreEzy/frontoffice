package role

import (
	"front-office/common/constant"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{Svc: service}
}

type controller struct {
	Svc Service
}

type Controller interface {
	GetRoleById(c *fiber.Ctx) error
	GetAllRoles(c *fiber.Ctx) error
}

func (ctrl *controller) GetRoleById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := ctrl.Svc.GetRoleById(id)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result == nil || result.Data.RoleId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a role by Id",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetAllRoles(c *fiber.Ctx) error {
	result, err := ctrl.Svc.GetAllRoles()
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get list of roles",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
