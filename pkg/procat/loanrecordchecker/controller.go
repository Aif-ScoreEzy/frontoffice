package loanrecordchecker

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
	LoanRecordChecker(c *fiber.Ctx) error
}

func (ctrl *controller) LoanRecordChecker(c *fiber.Ctx) error {
	req := c.Locals("request").(*LoanRecordCheckerRequest)
	apiKey, _ := c.Locals("apiKey").(string)

	res, err := ctrl.Svc.CallLoanRecordChecker(req, apiKey)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process loan record checker"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	result := LoanRecordCheckerResponse{
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
