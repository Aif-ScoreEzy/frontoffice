package grading

import (
	"fmt"
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
	GetGradings(c *fiber.Ctx) error
}

func (ctrl *controller) GetGradings(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	result, err := ctrl.Svc.GetGradings(companyId)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	res := helper.ResponseSuccess(
		"succeed to get gradings",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(res)
}
