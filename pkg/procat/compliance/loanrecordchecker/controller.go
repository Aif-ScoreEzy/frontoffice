package loanrecordchecker

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/apperror"

	"github.com/gofiber/fiber/v2"
)

func NewController(
	svc Service,
) Controller {
	return &controller{svc}
}

type controller struct {
	svc Service
}

type Controller interface {
	SingleSearch(c *fiber.Ctx) error
	BulkSearch(c *fiber.Ctx) error
}

func (ctrl *controller) SingleSearch(c *fiber.Ctx) error {
	reqBody := c.Locals("request").(*loanRecordCheckerRequest)
	apiKey := fmt.Sprintf("%v", c.Locals("apiKey"))
	memberIdStr := fmt.Sprintf("%v", c.Locals("userId"))
	companyIdStr := fmt.Sprintf("%v", c.Locals("companyId"))

	result, err := ctrl.svc.LoanRecordChecker(apiKey, memberIdStr, companyIdStr, reqBody)
	if err != nil {
		return err
	}

	return c.Status(result.StatusCode).JSON(result)
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
	apiKey := fmt.Sprintf("%v", c.Locals("apiKey"))

	memberId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidUserSession)
	}

	companyId, err := helper.InterfaceToUint(c.Locals("companyId"))
	if err != nil {
		return apperror.Unauthorized(constant.InvalidCompanySession)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	err = ctrl.svc.BulkLoanRecordChecker(apiKey, memberId, companyId, file)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"success",
		nil,
	))
}
