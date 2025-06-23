package taxcompliancestatus

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
	TaxComplianceStatus(c *fiber.Ctx) error
}

func (ctrl *controller) TaxComplianceStatus(c *fiber.Ctx) error {
	req := c.Locals("request").(*taxComplianceStatusRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	productRes, err := ctrl.productSvc.GetProductBySlug(constant.SlugTaxComplianceStatus)
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
	taxComplianceRes, err := ctrl.svc.CallTaxCompliance(apiKey, jobIdStr, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if taxComplianceRes.StatusCode >= fiber.StatusBadRequest {
		_, resp := helper.GetError("failed to process tax compliance status")

		return c.Status(taxComplianceRes.StatusCode).JSON(resp)
	}

	if err := ctrl.transactionSvc.UpdateLogProCat(taxComplianceRes.TransactionId, &transaction.UpdateTransRequest{
		Success: helper.BoolPtr(true),
	}); err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	logTransRes, err := ctrl.transactionSvc.GetLogTransSuccessCount(jobIdStr)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.logSvc.UpdateJobAPI(jobIdStr, &log.UpdateJobRequest{
		SuccessCount: &logTransRes.SuccessCount,
		Status:       helper.StringPtr(constant.JobStatusDone),
		EndAt:        helper.TimePtr(time.Now()),
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	return c.Status(taxComplianceRes.StatusCode).JSON(taxComplianceRes)
}
