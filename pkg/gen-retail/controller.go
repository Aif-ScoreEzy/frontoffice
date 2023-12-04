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
			Data:         nil,
		}

		return c.Status(genRetailResponse.StatusCode).JSON(dataReturn)
	}

	dataReturn := GenRetailV3DataClient{
		TransactionID:        genRetailResponse.Data.TransactionID,
		Name:                 genRetailResponse.Data.Name,
		IDCardNo:             genRetailResponse.Data.IDCardNo,
		PhoneNo:              genRetailResponse.Data.PhoneNo,
		LoanNo:               genRetailResponse.Data.LoanNo,
		ProbabilityToDefault: genRetailResponse.Data.ProbabilityToDefault,
		Grade:                genRetailResponse.Data.Grade,
		Date:                 genRetailResponse.Data.Date,
	}

	resp := GenRetailV3ClientReturnSuccess{
		Message: genRetailResponse.Message,
		Success: true,
		Data:    dataReturn,
	}

	return c.Status(genRetailResponse.StatusCode).JSON(resp)
}
