package livestatus

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
	UploadCSV(c *fiber.Ctx) error
}

func (ctrl *controller) UploadCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	csvData, totalData, err := helper.ParseCSVFile(file)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var fifRequest FIFRequest
	var FifRequests FIFRequests
	for i, request := range csvData {
		if i == 0 {
			continue
		}

		fifRequest.PhoneNumber = request[0]

		FifRequests.PhoneNumbers = append(FifRequests.PhoneNumbers, fifRequest)
	}

	err = ctrl.Svc.CreateJob(&FifRequests, totalData)
	if err != nil {
		return err
	}

	resp := helper.ResponseSuccess(
		"succeed to upload data",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
