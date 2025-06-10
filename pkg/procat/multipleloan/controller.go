package multipleloan

import (
	"front-office/helper"
	"strconv"

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
	MultipleLoan90Days(c *fiber.Ctx) error
}

func (ctrl *controller) MultipleLoan7Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*MultipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	res, err := ctrl.Svc.CallMultipleLoan7Days(req, apiKey, memberIdStr, companyIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan record checker"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) MultipleLoan30Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*MultipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	res, err := ctrl.Svc.CallMultipleLoan30Days(req, apiKey, memberIdStr, companyIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan record checker"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) MultipleLoan90Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*MultipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	res, err := ctrl.Svc.CallMultipleLoan90Days(req, apiKey, memberIdStr, companyIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan record checker"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
