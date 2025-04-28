package multipleloan

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func NewController(svc Service) Controller {
	return &controller{Svc: svc}
}

type controller struct {
	Svc Service
}
type Controller interface {
	MultipleLoan7Days(c *fiber.Ctx) error
	MultipleLoan30Days(c *fiber.Ctx) error
}

func (ctrl *controller) MultipleLoan7Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*MultipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)

	res, err := ctrl.Svc.CallMultipleLoan7Days(req, apiKey)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan record checker"
		}

		statusCode, resp := helper.GetError(msg)

		return c.Status(statusCode).JSON(resp)
	}

	result := MultipleLoanResponse{
		Data:            res.Data,
		PricingStrategy: res.PricingStrategy,
		TransactionID:   res.TransactionId,
		Datetime:        res.DateTime,
	}

	resp := helper.ResponseSuccess(
		"success",
		result,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) MultipleLoan30Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*MultipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)

	res, err := ctrl.Svc.CallMultipleLoan30Days(req, apiKey)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan record checker"
		}

		statusCode, resp := helper.GetError(msg)

		return c.Status(statusCode).JSON(resp)
	}

	result := MultipleLoanResponse{
		Data:            res.Data,
		PricingStrategy: res.PricingStrategy,
		TransactionID:   res.TransactionId,
		Datetime:        res.DateTime,
	}

	resp := helper.ResponseSuccess(
		"success",
		result,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
