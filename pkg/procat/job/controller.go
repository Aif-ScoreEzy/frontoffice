package job

import (
	"bytes"
	"errors"
	"fmt"
	"front-office/common/constant"
	"front-office/internal/apperror"
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
	GetJob(c *fiber.Ctx) error
	GetJobDetail(c *fiber.Ctx) error
	GetJobDetails(c *fiber.Ctx) error
	ExportJobDetail(c *fiber.Ctx) error
}

func (ctrl *controller) GetJob(c *fiber.Ctx) error {
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
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

func (ctrl *controller) GetJobDetail(c *fiber.Ctx) error {
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	filter := &logFilter{
		MemberId:    fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId:   fmt.Sprintf("%v", c.Locals("companyId")),
		Page:        c.Query("page", ""),
		Size:        c.Query("size", ""),
		ProductSlug: productSlug,
		JobId:       c.Params("job_id"),
	}

	result, err := ctrl.Svc.GetProCatJobDetail(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *controller) GetJobDetails(c *fiber.Ctx) error {
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate == "" {
		endDate = startDate
	}

	filter := &logFilter{
		MemberId:    fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId:   fmt.Sprintf("%v", c.Locals("companyId")),
		Page:        c.Query("page", ""),
		Size:        c.Query("size", ""),
		ProductSlug: productSlug,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	result, err := ctrl.Svc.GetProCatJobDetails(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *controller) ExportJobDetail(c *fiber.Ctx) error {
	memberID := c.Locals("userId").(uint)
	companyID := c.Locals("companyId").(uint)
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	filter := &logFilter{
		MemberId:    strconv.FormatUint(uint64(memberID), 10),
		CompanyId:   strconv.FormatUint(uint64(companyID), 10),
		ProductSlug: productSlug,
		JobId:       c.Params("job_id"),
	}

	var buf bytes.Buffer
	filename, err := ctrl.Svc.ExportJobDetailToCSV(filter, &buf)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

var productSlugMap = map[string]string{
	"loan-record-checker":     constant.SlugLoanRecordChecker,
	"7d-multiple-loan":        constant.SlugMultipleLoan7Days,
	"30d-multiple-loan":       constant.SlugMultipleLoan30Days,
	"90d-multiple-loan":       constant.SlugMultipleLoan90Days,
	"tax-compliance-status":   constant.SlugTaxComplianceStatus,
	"tax-score":               constant.SlugTaxScore,
	"tax-verification-detail": constant.SlugTaxVerificationDetail,
}

func mapProductSlug(slug string) (string, error) {
	if mapped, ok := productSlugMap[slug]; ok {
		return mapped, nil
	}

	return "", errors.New("unsupported product slug")
}
