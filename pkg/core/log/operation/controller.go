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
	GetList(c *fiber.Ctx) error
	GetListByRange(c *fiber.Ctx) error
}

func (ctrl *controller) GetList(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	role := c.Query("role")
	event := c.Query("event")

	filter := &LogOperationFilter{
		CompanyId: companyId,
		Role:      role,
		Event:     event,
	}

	result, err := ctrl.Svc.GetLogOperations(filter)
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

func (ctrl *controller) GetListByRange(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	startDate := c.Query("start_date")
	endDate := c.Query(("end_date"))

	filter := &LogRangeFilter{
		CompanyId: companyId,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := ctrl.Svc.GetByRange(filter)
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
