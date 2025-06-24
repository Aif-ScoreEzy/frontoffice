package loanrecordchecker

import (
	"fmt"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/procat/log"

	"github.com/gofiber/fiber/v2"
)

func NewController(
	svc Service,
	productSvc product.Service,
	logSvc log.Service,
	transactionSvc transaction.Service,
) Controller {
	return &controller{svc, productSvc, logSvc, transactionSvc}
}

type controller struct {
	svc            Service
	productSvc     product.Service
	logSvc         log.Service
	transactionSvc transaction.Service
}

type Controller interface {
	LoanRecordChecker(c *fiber.Ctx) error
}

func (ctrl *controller) LoanRecordChecker(c *fiber.Ctx) error {
	req := c.Locals("request").(*LoanRecordCheckerRequest)
	apiKey := fmt.Sprintf("%v", c.Locals("apiKey"))
	memberIdStr := fmt.Sprintf("%v", c.Locals("userId"))
	companyIdStr := fmt.Sprintf("%v", c.Locals("companyId"))

	result, err := ctrl.svc.LoanRecordChecker(apiKey, memberIdStr, companyIdStr, req)
	if err != nil {
		return err
	}

	return c.Status(result.StatusCode).JSON(result)
}
