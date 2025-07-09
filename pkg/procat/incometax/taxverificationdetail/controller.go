package taxverificationdetail

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"

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
	TaxVerificationDetail(c *fiber.Ctx) error
}

func (ctrl *controller) TaxVerificationDetail(c *fiber.Ctx) error {
	reqBody := c.Locals(constant.Request).(*taxVerificationRequest)
	apiKey, _ := c.Locals(constant.APIKey).(string)
	memberId := fmt.Sprintf("%v", c.Locals(constant.UserId))
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	result, err := ctrl.svc.CallTaxVerification(apiKey, memberId, companyId, reqBody)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	return c.Status(result.StatusCode).JSON(result)
}
