package log

import (
	"errors"
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

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.ResponseFailed(err.Error()))
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
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *controller) GetProCatJobDetail(c *fiber.Ctx) error {
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.ResponseFailed(err.Error()))
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
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func mapProductSlug(slug string) (string, error) {
	switch slug {
	case "loan-record-checker":
		return constant.SlugLoanRecordChecker, nil
	case "7d-multiple-loan":
		return constant.SlugMultipleLoan7Days, nil
	case "30d-multiple-loan":
		return constant.SlugMultipleLoan30Days, nil
	case "90d-multiple-loan":
		return constant.SlugMultipleLoan90Days, nil
	case "tax-compliance-status":
		return constant.SlugTaxComplianceStatus, nil
	case "tax-score":
		return constant.SlugTaxScore, nil
	case "tax-verification-detail":
		return constant.SlugTaxVerificationDetail, nil
	default:
		return "", errors.New("unsupported product slug")
	}
}
