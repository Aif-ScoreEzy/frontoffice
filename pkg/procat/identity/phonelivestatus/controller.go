package phonelivestatus

import (
	"bytes"
	"fmt"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/member"

	"github.com/gofiber/fiber/v2"
)

func NewController(svc Service, memberSvc member.Service) Controller {
	return &controller{svc, memberSvc}
}

type controller struct {
	svc       Service
	memberSvc member.Service
}

type Controller interface {
	GetJobs(c *fiber.Ctx) error
	GetJobDetails(c *fiber.Ctx) error
	ExportJobDetails(c *fiber.Ctx) error
	GetJobsSummary(c *fiber.Ctx) error
	ExportJobsSummary(c *fiber.Ctx) error
	SingleSearch(c *fiber.Ctx) error
	BulkSearch(c *fiber.Ctx) error
}

func (ctrl *controller) GetJobs(c *fiber.Ctx) error {
	filter := &phoneLiveStatusFilter{
		Page:      c.Query("page", "1"),
		Size:      c.Query("size", "10"),
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	jobs, err := ctrl.svc.GetPhoneLiveStatusJob(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeeded to get phone live status jobs",
		jobs,
	))
}

func (ctrl *controller) GetJobDetails(c *fiber.Ctx) error {
	filter := &phoneLiveStatusFilter{
		Page:      c.Query("page", "1"),
		Size:      c.Query("size", "10"),
		JobId:     c.Params("id"),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	if filter.JobId == "" {
		return apperror.BadRequest("missing job ID")
	}

	jobDetail, err := ctrl.svc.GetPhoneLiveStatusDetailsSummary(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeeded to get phone live status job details",
		jobDetail,
	))
}

func (ctrl *controller) GetJobsSummary(c *fiber.Ctx) error {
	filter := &phoneLiveStatusFilter{
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	if filter.StartDate == "" || filter.EndDate == "" {
		return apperror.BadRequest("start_date and end_date are required")
	}

	jobsSummary, err := ctrl.svc.GetJobsSummary(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeeded to get phone live status jobs summary",
		jobsSummary,
	))
}

func (ctrl *controller) ExportJobsSummary(c *fiber.Ctx) error {
	filter := &phoneLiveStatusFilter{
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	if filter.StartDate == "" || filter.EndDate == "" {
		return apperror.BadRequest("start_date and end_date are required")
	}

	var buf bytes.Buffer
	filename, err := ctrl.svc.ExportJobsSummary(filter, &buf)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

func (ctrl *controller) ExportJobDetails(c *fiber.Ctx) error {
	filter := &phoneLiveStatusFilter{
		JobId:     c.Params("id"),
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	var buf bytes.Buffer

	filename, err := ctrl.svc.ExportJobDetails(filter, &buf)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

func (ctrl *controller) SingleSearch(c *fiber.Ctx) error {
	reqBody := c.Locals("request").(*phoneLiveStatusRequest)

	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	err := ctrl.svc.ProcessPhoneLiveStatus(memberId, companyId, reqBody)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"success",
		nil,
	))
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	file, err := c.FormFile("file")
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if err := helper.ValidateUploadedFile(file, 30*1024*1024, []string{".csv"}); err != nil {
		_, resp := helper.GetError(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	err = ctrl.svc.BulkProcessPhoneLiveStatus(memberId, companyId, file)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"success",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
