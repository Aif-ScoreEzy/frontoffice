package member

import (
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
}

func (ctrl *controller) GetBy(c *fiber.Ctx) error {
	id := c.Query("id")
	email := c.Query("email")
	username := c.Query("username")
	key := c.Query("key")

	result, err := ctrl.Svc.GetBy(&FindUserQuery{
		Id:       id,
		Email:    email,
		Username: username,
		Key:      key,
	})
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user by Id",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
