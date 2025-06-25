package loanrecordchecker

import (
	"fmt"

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
	LoanRecordChecker(c *fiber.Ctx) error
}

func (ctrl *controller) LoanRecordChecker(c *fiber.Ctx) error {
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
