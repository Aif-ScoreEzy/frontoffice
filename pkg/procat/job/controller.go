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
	GetJobDetails(c *fiber.Ctx) error
	ExportJobDetails(c *fiber.Ctx) error
	GetJobDetailsByDateRange(c *fiber.Ctx) error
	ExportJobDetailsByDateRange(c *fiber.Ctx) error
}

func (ctrl *controller) GetJob(c *fiber.Ctx) error {
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	filter := &logFilter{
		Page:        c.Query(constant.Page, "1"),
		Size:        c.Query(constant.Size, "10"),
		StartDate:   c.Query(constant.StartDate, ""),
		EndDate:     c.Query(constant.EndDate, ""),
		ProductSlug: productSlug,
		MemberId:    fmt.Sprintf("%v", c.Locals(constant.UserId)),
		CompanyId:   fmt.Sprintf("%v", c.Locals(constant.CompanyId)),
		TierLevel:   fmt.Sprintf("%v", c.Locals(constant.RoleId)),
	}

	result, err := ctrl.Svc.GetJob(filter)
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

	filter := &logFilter{
		MemberId:    fmt.Sprintf("%v", c.Locals(constant.UserId)),
		CompanyId:   fmt.Sprintf("%v", c.Locals(constant.CompanyId)),
		Page:        c.Query(constant.Page, ""),
		Size:        c.Query(constant.Size, ""),
		ProductSlug: productSlug,
		JobId:       c.Params("job_id"),
	}

	result, err := ctrl.Svc.GetJobDetails(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *controller) GetJobDetailsByDateRange(c *fiber.Ctx) error {
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	startDate := c.Query(constant.StartDate)
	endDate := c.Query(constant.EndDate)
	if startDate != "" && endDate == "" {
		endDate = startDate
	}

	filter := &logFilter{
		MemberId:    fmt.Sprintf("%v", c.Locals(constant.UserId)),
		CompanyId:   fmt.Sprintf("%v", c.Locals(constant.CompanyId)),
		Page:        c.Query(constant.Page, "1"),
		Size:        c.Query(constant.Size, "10"),
		ProductSlug: productSlug,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	result, err := ctrl.Svc.GetJobDetailsByDateRange(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *controller) ExportJobDetails(c *fiber.Ctx) error {
	memberId := c.Locals(constant.UserId).(uint)
	companyId := c.Locals(constant.CompanyId).(uint)
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	filter := &logFilter{
		MemberId:    strconv.FormatUint(uint64(memberId), 10),
		CompanyId:   strconv.FormatUint(uint64(companyId), 10),
		ProductSlug: productSlug,
		JobId:       c.Params("job_id"),
	}

	var buf bytes.Buffer
	filename, err := ctrl.Svc.ExportJobDetails(filter, &buf)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

func (ctrl *controller) ExportJobDetailsByDateRange(c *fiber.Ctx) error {
	memberId := c.Locals(constant.UserId).(uint)
	companyId := c.Locals(constant.CompanyId).(uint)
	slug := c.Params("product_slug")

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	startDate := c.Query(constant.StartDate)
	endDate := c.Query(constant.EndDate)
	if startDate != "" && endDate == "" {
		endDate = startDate
	}

	filter := &logFilter{
		MemberId:    strconv.FormatUint(uint64(memberId), 10),
		CompanyId:   strconv.FormatUint(uint64(companyId), 10),
		ProductSlug: productSlug,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	var buf bytes.Buffer
	filename, err := ctrl.Svc.ExportJobDetailsByDateRange(filter, &buf)
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
