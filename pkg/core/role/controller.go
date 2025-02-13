package role

import (
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/permission"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service, svcPermission permission.Service) Controller {
	return &controller{Svc: service, SvcPermission: svcPermission}
}

type controller struct {
	Svc           Service
	SvcPermission permission.Service
}

type Controller interface {
	GetRoleById(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	GetAllRoles(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
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

func (ctrl *controller) CreateRole(c *fiber.Ctx) error {
	request := c.Locals("request").(*CreateRoleRequest)

	_, err := ctrl.Svc.GetRoleByNameSvc(request.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	for _, permissionData := range request.Permissions {
		_, err := ctrl.SvcPermission.IsPermissionExistSvc(permissionData.Id)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}
	}

	role, err := ctrl.Svc.CreateRoleSvc(request)
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

func (ctrl *controller) UpdateRole(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateRoleRequest)
	id := c.Params("id")

	_, err := ctrl.Svc.GetRoleById(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if req.Name != "" {
		_, err = ctrl.Svc.GetRoleByNameSvc(req.Name)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}
	}

	role, err := ctrl.Svc.UpdateRoleByIdSvc(req, id)
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

func (ctrl *controller) DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := ctrl.Svc.GetRoleById(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if err := ctrl.Svc.DeleteRoleByIdSvc(id); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to delete a role",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
