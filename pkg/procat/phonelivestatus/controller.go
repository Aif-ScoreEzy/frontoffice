package phonelivestatus

import (
	"bytes"
	"fmt"
	"front-office/helper"
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
	filter := &PhoneLiveStatusFilter{
		Page:      c.Query("page", "1"),
		Size:      c.Query("size", "10"),
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	result, err := ctrl.svc.GetPhoneLiveStatusJob(filter)
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

func (ctrl *controller) GetJobDetails(c *fiber.Ctx) error {
	filter := &PhoneLiveStatusFilter{
		Page:      c.Query("page", "1"),
		Size:      c.Query("size", "10"),
		JobId:     c.Params("id"),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	result, err := ctrl.svc.GetPhoneLiveStatusDetailsSummary(filter)
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

func (ctrl *controller) GetJobsSummary(c *fiber.Ctx) error {
	filter := &PhoneLiveStatusFilter{
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	result, err := ctrl.svc.GetJobsSummary(filter)
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

func (ctrl *controller) ExportJobsSummary(c *fiber.Ctx) error {
	filter := &PhoneLiveStatusFilter{
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	result, err := ctrl.svc.GetPhoneLiveStatusDetailsByRangeDate(filter)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	var buf bytes.Buffer

	filename, err := ctrl.svc.ExportJobsSummary(result.Data, filter, &buf)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

func (ctrl *controller) ExportJobDetails(c *fiber.Ctx) error {
	filter := &PhoneLiveStatusFilter{
		JobId:     c.Params("id"),
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		MemberId:  fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId: fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel: fmt.Sprintf("%v", c.Locals("roleId")),
	}

	// Get JobDetails By JobId
	result, err := ctrl.svc.GetAllPhoneLiveStatusDetails(filter)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var buf bytes.Buffer

	filename, err := ctrl.svc.ExportJobsSummary(result.Data, filter, &buf)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

func (ctrl *controller) SingleSearch(c *fiber.Ctx) error {
	req := c.Locals("request").(*PhoneLiveStatusRequest)
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	err := ctrl.svc.ProcessPhoneLiveStatus(memberId, companyId, req)
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
