package loanrecordchecker

import (
	"fmt"
	"front-office/common/constant"
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
	LoanRecordChecker(c *fiber.Ctx) error
	GetLoanRecordCheckerJob(c *fiber.Ctx) error
}

func (ctrl *controller) LoanRecordChecker(c *fiber.Ctx) error {
	req := c.Locals("request").(*LoanRecordCheckerRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	res, err := ctrl.Svc.LoanRecordChecker(req, apiKey, memberIdStr, companyIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if !res.Success {
		msg := res.Message
		if msg == "" {
			msg = "failed to process loan record checker"
		}

		resp := helper.ResponseFailed(
			msg,
		)

		return c.Status(res.StatusCode).JSON(resp)
	}

	result := LoanRecordCheckerResponse{
		Data:            res.Data,
		PricingStrategy: res.PricingStrategy,
		TransactionID:   res.TransactionId,
		Datetime:        res.DateTime,
	}

	resp := helper.ResponseSuccess(
		"success",
		result,
	)

	return c.Status(res.StatusCode).JSON(resp)
}

func (ctrl *controller) GetLoanRecordCheckerJob(c *fiber.Ctx) error {
	filter := &loanRecordCheckerFilter{
		Page:        c.Query("page", "1"),
		Size:        c.Query("size", "10"),
		StartDate:   c.Query("start_date", ""),
		EndDate:     c.Query("end_date", ""),
		ProductSlug: constant.SlugLoanRecordChecker,
		MemberId:    fmt.Sprintf("%v", c.Locals("userId")),
		CompanyId:   fmt.Sprintf("%v", c.Locals("companyId")),
		TierLevel:   fmt.Sprintf("%v", c.Locals("roleId")),
	}

	result, err := ctrl.Svc.GetLoanRecordCheckerJob(filter)
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
