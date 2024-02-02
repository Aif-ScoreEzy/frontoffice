package log

import (
	"front-office/common/model"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func GetAllLogTransByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	companyID := c.Query("company_id")

	result, statusCode, errRequest := GetAllLogTransByDateSvc(companyID, date)
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
