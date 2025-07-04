package transaction

import (
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/apperror"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *controller) GetLogScoreezy(c *fiber.Ctx) error {
	logs, err := ctrl.svc.GetScoreezyLogs()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		constant.SucceedGetLogTrans,
		logs,
	))
}

func (ctrl *controller) GetLogScoreezyByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	companyId := c.Query("company_id")

	if date == "" {
		return apperror.BadRequest("date are required")
	}

	logs, err := ctrl.svc.GetScoreezyLogsByDate(companyId, date)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		constant.SucceedGetLogTrans,
		logs,
	))
}

func (ctrl *controller) GetLogScoreezyByRangeDate(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	companyId := c.Query("company_id")

	if startDate == "" || endDate == "" {
		return apperror.BadRequest("start_date and end_date  are required")
	}

	logs, err := ctrl.svc.GetScoreezyLogsByRangeDate(startDate, endDate, companyId, page)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		constant.SucceedGetLogTrans,
		logs,
	))
}

func (ctrl *controller) GetLogScoreezyByMonth(c *fiber.Ctx) error {
	companyId := c.Query("company_id")
	month := c.Query("month")

	if month == "" {
		return apperror.BadRequest("month are required")
	}

	logs, err := ctrl.svc.GetScoreezyLogsByMonth(companyId, month)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		constant.SucceedGetLogTrans,
		logs,
	))
}
