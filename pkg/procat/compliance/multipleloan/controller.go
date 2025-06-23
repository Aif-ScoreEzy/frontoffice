package multipleloan

import (
	"errors"
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
	MultipleLoan(c *fiber.Ctx) error
}

func (ctrl *controller) MultipleLoan(c *fiber.Ctx) error {
	req := c.Locals("request").(*multipleLoanRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId, _ := c.Locals("userId").(uint)
	companyId, _ := c.Locals("companyId").(uint)
	slug := c.Params("product_slug")

	memberIdStr := strconv.FormatUint(uint64(memberId), 10)
	companyIdStr := strconv.FormatUint(uint64(companyId), 10)

	productSlug, err := mapProductSlug(slug)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.ResponseFailed(err.Error()))
	}

	productRes, err := ctrl.productSvc.GetProductBySlug(productSlug)
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

	multipleLoanRes, err := ctrl.svc.MultipleLoan(apiKey, jobIdStr, productSlug, memberIdStr, companyIdStr, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}
	if multipleLoanRes.StatusCode >= fiber.StatusBadRequest {
		_, resp := helper.GetError(multipleLoanRes.Message)

		return c.Status(multipleLoanRes.StatusCode).JSON(resp)
	}

	if err := ctrl.transactionSvc.UpdateLogProCat(multipleLoanRes.TransactionId, &transaction.UpdateTransRequest{
		Success: helper.BoolPtr(true),
	}); err != nil {
		return err
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

	return c.Status(fiber.StatusOK).JSON(multipleLoanRes)
}

var productSlugMap = map[string]string{
	"7d-multiple-loan":  constant.SlugMultipleLoan7Days,
	"30d-multiple-loan": constant.SlugMultipleLoan30Days,
	"90d-multiple-loan": constant.SlugMultipleLoan90Days,
}

func mapProductSlug(slug string) (string, error) {
	if val, ok := productSlugMap[slug]; ok {
		return val, nil
	}

	return "", errors.New("unsupported product slug")
}
