package log

import (
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
	GetTransactionLogs(c *fiber.Ctx) error
	GetTransactionLogsByDate(c *fiber.Ctx) error
	GetTransactionLogsByRangeDate(c *fiber.Ctx) error
	GetTransactionLogsByMonth(c *fiber.Ctx) error
	GetTransactionLogsByName(c *fiber.Ctx) error
}

func (ctrl *controller) GetTransactionLogs(c *fiber.Ctx) error {
	result, statusCode, errRequest := ctrl.Svc.GetTransactionLogsSvc()
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}

func (ctrl *controller) GetTransactionLogsByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	companyID := c.Query("company_id")

	result, statusCode, errRequest := ctrl.Svc.GetTransactionLogsByDateSvc(companyID, date)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}

func (ctrl *controller) GetTransactionLogsByRangeDate(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	companyID := c.Query("company_id")
	page := c.Query("page", "1")

	result, statusCode, errRequest := ctrl.Svc.GetTransactionLogsByRangeDateSvc(startDate, endDate, companyID, page)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}

func (ctrl *controller) GetTransactionLogsByMonth(c *fiber.Ctx) error {
	companyID := c.Query("company_id")
	month := c.Query("month")

	result, statusCode, errRequest := ctrl.Svc.GetTransactionLogsByMonthSvc(companyID, month)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}

func (ctrl *controller) GetTransactionLogsByName(c *fiber.Ctx) error {
	companyID := c.Query("company_id")
	name := c.Query("name")

	result, statusCode, errRequest := ctrl.Svc.GetTransactionLogsByNameSvc(companyID, name)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}
