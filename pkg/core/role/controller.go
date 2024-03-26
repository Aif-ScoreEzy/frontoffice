package role

import (
	"front-office/helper"
	"front-office/pkg/core/permission"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{Svc: service}
}

type controller struct {
	Svc           Service
	SvcPermission permission.Service
}

type Controller interface {
	CreateRole(c *fiber.Ctx) error
	GetAllRoles(c *fiber.Ctx) error
	GetRoleByID(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
}

func (ctrl *controller) CreateRole(c *fiber.Ctx) error {
	request := c.Locals("request").(*CreateRoleRequest)

	_, err := ctrl.Svc.GetRoleByNameSvc(request.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	for _, permissionData := range request.Permissions {
		_, err := ctrl.SvcPermission.IsPermissionExistSvc(permissionData.ID)
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
	roles, err := ctrl.Svc.GetAllRolesSvc()
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

func (ctrl *controller) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	role, err := ctrl.Svc.FindRoleByIDSvc(id)
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

func (ctrl *controller) UpdateRole(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateRoleRequest)
	id := c.Params("id")

	_, err := ctrl.Svc.FindRoleByIDSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	_, err = ctrl.Svc.GetRoleByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	role, err := ctrl.Svc.UpdateRoleByIDSvc(req, id)
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

	_, err := ctrl.Svc.FindRoleByIDSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if err := ctrl.Svc.DeleteRoleByIDSvc(id); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to delete a role",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
