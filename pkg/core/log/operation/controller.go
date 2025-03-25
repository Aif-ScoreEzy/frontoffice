package operation

import (
	"fmt"
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
	GetLogOperationsByCompany(c *fiber.Ctx) error
}

func (ctrl *controller) GetLogOperationsByCompany(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	role := c.Query("role")
	event := c.Query("event")

	filter := &GetLogOperationFilter{
		Role:  role,
		Event: event,
	}

	result, err := ctrl.Svc.GetLogOperations(companyId, filter)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(res)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log operation",
		result,
	)

	return c.Status(responseBody.StatusCode).JSON(responseBody)
}
