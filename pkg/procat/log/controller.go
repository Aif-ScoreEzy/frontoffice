package log

import (
	"fmt"
	"front-office/common/constant"
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
	GetProCatJob(c *fiber.Ctx) error
	GetProCatJobDetail(c *fiber.Ctx) error
}

func (ctrl *controller) GetProCatJob(c *fiber.Ctx) error {
	slug := c.Params("product_slug")
	var productSlug string

	switch slug {
	case "loan-record-checker":
		productSlug = constant.SlugLoanRecordChecker
	case "7d-multiple-loan":
		productSlug = constant.SlugMultipleLoan7Days
	case "30d-multiple-loan":
		productSlug = constant.SlugMultipleLoan30Days
	case "90d-multiple-loan":
		productSlug = constant.SlugMultipleLoan90Days
	case "tax-compliance-status":
		productSlug = constant.SlugTaxComplianceStatus
	case "tax-score":
		productSlug = constant.SlugTaxScore
	default:
		return c.Status(fiber.StatusNotFound).JSON(helper.ResponseFailed("Unsupported product slug"))
	}

	filter := &logFilter{
		Page:        c.Query("page", "1"),
		Size:        c.Query("size", "10"),
		StartDate:   c.Query("start_date", ""),
		EndDate:     c.Query("end_date", ""),
		ProductSlug: productSlug,
		MemberId:    fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId:   fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel:   fmt.Sprintf("%v", c.Locals("roleId")),
	}

	result, err := ctrl.Svc.GetProCatJob(filter)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if result.StatusCode != fiber.StatusOK {
		_, resp := helper.GetError(result.Message)

		return c.Status(result.StatusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *controller) GetProCatJobDetail(c *fiber.Ctx) error {
	slug := c.Params("product_slug")
	var productSlug string

	switch slug {
	case "loan-record-checker":
		productSlug = constant.SlugLoanRecordChecker
	case "7d-multiple-loan":
		productSlug = constant.SlugMultipleLoan7Days
	case "30d-multiple-loan":
		productSlug = constant.SlugMultipleLoan30Days
	case "90d-multiple-loan":
		productSlug = constant.SlugMultipleLoan90Days
	case "tax-compliance-status":
		productSlug = constant.SlugTaxComplianceStatus
	case "tax-score":
		productSlug = constant.SlugTaxScore
	default:
		return c.Status(fiber.StatusNotFound).JSON(helper.ResponseFailed("Unsupported product slug"))
	}

	filter := &logFilter{
		MemberId:    fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId:   fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel:   fmt.Sprintf("%v", c.Locals("roleId")),
		ProductSlug: productSlug,
		JobId:       c.Params("job_id"),
	}

	result, err := ctrl.Svc.GetProCatJobDetail(filter)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if result.StatusCode != fiber.StatusOK {
		_, resp := helper.GetError(result.Message)

		return c.Status(result.StatusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
