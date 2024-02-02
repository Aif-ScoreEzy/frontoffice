package log

import (
	"front-office/common/model"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func GetTransactionLogsByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	companyID := c.Query("company_id")

	result, statusCode, errRequest := GetTransactionLogsByDateSvc(companyID, date)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := model.AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}

func GetTransactionLogsByRangeDate(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	companyID := c.Query("company_id")
	page := c.Query("page", "1")

	result, statusCode, errRequest := GetTransactionLogsByRangeDateSvc(startDate, endDate, companyID, page)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := model.AifResponse{
		Data: result.Data,
		Meta: result.Meta,
	}

	return c.Status(statusCode).JSON(resp)
}
