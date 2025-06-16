package multipleloan

import (
	"front-office/common/constant"
	"front-office/common/model"
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
	MultipleLoan(c *fiber.Ctx) error
	MultipleLoan7Days(c *fiber.Ctx) error
	MultipleLoan30Days(c *fiber.Ctx) error
	MultipleLoan90Days(c *fiber.Ctx) error
}

func (ctrl *controller) MultipleLoan(c *fiber.Ctx) error {
	req := c.Locals("request").(*multipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	slug := c.Params("product_slug")
	var res *model.ProCatAPIResponse[dataMultipleLoanResponse]
	var err error

	switch slug {
	case constant.SlugMultipleLoan7Days:
		res, err = ctrl.Svc.CallMultipleLoan7Days(req, apiKey, memberIdStr, companyIdStr)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())

			return c.Status(statusCode).JSON(resp)
		}

	case constant.SlugMultipleLoan30Days:
		res, err = ctrl.Svc.CallMultipleLoan30Days(req, apiKey, memberIdStr, companyIdStr)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())

			return c.Status(statusCode).JSON(resp)
		}

	case constant.SlugMultipleLoan90Days:
		res, err = ctrl.Svc.CallMultipleLoan90Days(req, apiKey, memberIdStr, companyIdStr)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())

			return c.Status(statusCode).JSON(resp)
		}

	default:
		return c.Status(fiber.StatusNotFound).JSON(helper.ResponseFailed("Unsupported product slug"))
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan request"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) MultipleLoan7Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*multipleLoanRequest)
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

	if res.StatusCode > fiber.StatusBadRequest {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(res.StatusCode).JSON(res)
}

func (ctrl *controller) MultipleLoan30Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*multipleLoanRequest)
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

	if res.StatusCode > fiber.StatusBadRequest {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(res.StatusCode).JSON(res)
}

func (ctrl *controller) MultipleLoan90Days(c *fiber.Ctx) error {
	req := c.Locals("request").(*multipleLoanRequest)
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

	if res.StatusCode > fiber.StatusBadRequest {
		msg := res.Message
		if msg == "" {
			msg = "failed to process multiple loan"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	return c.Status(res.StatusCode).JSON(res)
}
