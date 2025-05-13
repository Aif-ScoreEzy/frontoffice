package phonelivestatus

import (
	"fmt"
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
	GetJobs(c *fiber.Ctx) error
	SingleSearch(c *fiber.Ctx) error
	BulkSearch(c *fiber.Ctx) error
}

func (ctrl *controller) GetJobs(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	tierLevel := fmt.Sprintf("%v", c.Locals("roleId"))

	filter := &PhoneLiveStatusFilter{
		Page:      page,
		Size:      size,
		StartDate: startDate,
		EndDate:   endDate,
		MemberId:  memberId,
		CompanyId: companyId,
		TierLevel: tierLevel,
	}

	result, err := ctrl.Svc.GetPhoneLiveStatusJobAPI(filter)
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

func (ctrl *controller) SingleSearch(c *fiber.Ctx) error {
	req := c.Locals("request").(*PhoneLiveStatusRequest)
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	err := ctrl.Svc.ProcessPhoneLiveStatus(memberId, companyId, req)
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

	err = ctrl.Svc.BulkProcessPhoneLiveStatus(memberId, companyId, file)
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
