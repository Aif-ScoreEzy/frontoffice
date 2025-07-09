package permission

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
	CreatePermission(c *fiber.Ctx) error
	GetPermissionById(c *fiber.Ctx) error
	UpdatePermissionById(c *fiber.Ctx) error
	DeletePermissionById(c *fiber.Ctx) error
}

func (ctrl *controller) CreatePermission(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*PermissionRequest)

	_, err := ctrl.Svc.GetPermissionByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	permission, err := ctrl.Svc.CreatePermissionSvc(*req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to create a permission",
		permission,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (ctrl *controller) GetPermissionById(c *fiber.Ctx) error {
	id := c.Params("id")

	permission, err := ctrl.Svc.IsPermissionExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a permission by Id",
		permission,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UpdatePermissionById(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*PermissionRequest)
	id := c.Params("id")

	_, err := ctrl.Svc.IsPermissionExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	_, err = ctrl.Svc.GetPermissionByNameSvc(req.Name)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	permission, err := ctrl.Svc.UpdatePermissionByIdSvc(*req, id)
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

func (ctrl *controller) DeletePermissionById(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := ctrl.Svc.IsPermissionExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if err := ctrl.Svc.DeletePermissionByIdSvc(id); err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to delete a permission",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
