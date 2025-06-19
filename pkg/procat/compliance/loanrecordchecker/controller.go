package loanrecordchecker

import (
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/procat/log"
	"strconv"
	"time"

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
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	productRes, err := ctrl.productSvc.GetProductBySlug(constant.SlugLoanRecordChecker)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	jobRes, err := ctrl.logSvc.CreateProCatJob(&log.CreateJobRequest{
		ProductId: productRes.Data.ProductId,
		MemberId:  memberIdStr,
		CompanyId: companyIdStr,
		Total:     1,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	jobIdStr := strconv.FormatUint(uint64(jobRes.Data.JobId), 10)
	loanRecordRes, err := ctrl.svc.LoanRecordChecker(req, apiKey, jobIdStr, memberIdStr, companyIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if loanRecordRes.StatusCode > fiber.StatusBadRequest {
		_, resp := helper.GetError(loanRecordRes.Data.Status)

		return c.Status(loanRecordRes.StatusCode).JSON(resp)
	}

	_, err = ctrl.transactionSvc.UpdateLogProCat(loanRecordRes.TransactionId, &transaction.UpdateTransRequest{
		Success: helper.BoolPtr(true),
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	logTransRes, err := ctrl.transactionSvc.GetLogTransSuccessCount(jobIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.logSvc.UpdateJobAPI(jobIdStr, &log.UpdateJobRequest{
		SuccessCount: &logTransRes.Data.SuccessCount,
		Status:       helper.StringPtr(constant.JobStatusDone),
		EndAt:        helper.TimePtr(time.Now()),
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	return c.Status(loanRecordRes.StatusCode).JSON(loanRecordRes)
}
