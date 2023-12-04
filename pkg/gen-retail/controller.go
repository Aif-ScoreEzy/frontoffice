package genretail

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func RequestScore(c *fiber.Ctx) error {
	req := c.Locals("request").(*GenRetailRequest)
	apiKey := c.Get("X-API-KEY")

	genRetailResponse, errRequest := GenRetailV3(req, apiKey)
	if errRequest != nil {
		statusCode, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if genRetailResponse.StatusCode >= 400 {
		dataReturn := GenRetailV3ClientReturnError{
			Message:      genRetailResponse.Message,
			ErrorMessage: genRetailResponse.ErrorMessage,
			Data:         genRetailResponse.Data,
		}

		return c.Status(genRetailResponse.StatusCode).JSON(dataReturn)
	}

	resp := GenRetailV3ClientReturnSuccess{
		Message: genRetailResponse.Message,
		Success: true,
		Data:    genRetailResponse.Data,
	}

	return c.Status(genRetailResponse.StatusCode).JSON(resp)
}
