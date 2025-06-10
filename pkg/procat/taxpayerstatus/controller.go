package taxpayerstatus

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func NewController(svc Service) Controller {
	return &controller{svc}
}

type controller struct {
	svc Service
}

type Controller interface {
	TaxPayerStatus(c *fiber.Ctx) error
}

func (ctrl *controller) TaxPayerStatus(c *fiber.Ctx) error {
	req := c.Locals("request").(*taxPayerStatusRequest)
	apiKey, _ := c.Locals("apiKey").(string)

	res, err := ctrl.svc.CallTaxPayerStatus(apiKey, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if res.StatusCode >= 400 {
		_, resp := helper.GetError(res.Message)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(res.StatusCode).JSON(res)
}
