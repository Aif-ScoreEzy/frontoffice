package member

import (
	"front-office/common/constant"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{
		Svc: service,
	}
}

type controller struct {
	Svc Service
}

type Controller interface {
	GetBy(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetList(c *fiber.Ctx) error
	DeleteById(c *fiber.Ctx) error
}

func (ctrl *controller) GetBy(c *fiber.Ctx) error {
	email := c.Query("email")
	username := c.Query("username")
	key := c.Query("key")

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Email:    email,
		Username: username,
		Key:      key,
	})
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: id,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result == nil || !result.Success || result.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetList(c *fiber.Ctx) error {
	result, err := ctrl.Svc.GetMemberList()
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get member list",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: id,
	})

	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result == nil || !result.Success || result.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	result, err = ctrl.Svc.DeleteMemberById(id)
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to delete member",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
