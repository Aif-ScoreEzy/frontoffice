package taxverificationdetail

import (
	"fmt"
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
	TaxVerificationDetail(c *fiber.Ctx) error
}

func (ctrl *controller) TaxVerificationDetail(c *fiber.Ctx) error {
	req := c.Locals("request").(*taxVerificationRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	productRes, err := ctrl.productSvc.GetProductBySlug(constant.SlugTaxVerificationDetail)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	jobRes, err := ctrl.logSvc.CreateProCatJob(&log.CreateJobRequest{
		ProductId: productRes.Data.ProductId,
		MemberId:  memberId,
		CompanyId: companyId,
		Total:     1,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	jobIdStr := strconv.FormatUint(uint64(jobRes.Data.JobId), 10)
	taxVerificationRes, err := ctrl.svc.CallTaxVerification(apiKey, jobIdStr, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if taxVerificationRes.StatusCode >= fiber.StatusBadRequest {
		_, resp := helper.GetError("failed to process tax verification detail")

		return c.Status(taxVerificationRes.StatusCode).JSON(resp)
	}

	_, err = ctrl.transactionSvc.UpdateLogProCat(taxVerificationRes.TransactionId, &transaction.UpdateTransRequest{
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

	return c.Status(taxVerificationRes.StatusCode).JSON(taxVerificationRes)
}
