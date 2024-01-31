package log

import (
	"front-office/common/model"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func GetAllLogTrans(c *fiber.Ctx) error {
	result, statusCode, errRequest := GetAllLogTransSvc()
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
