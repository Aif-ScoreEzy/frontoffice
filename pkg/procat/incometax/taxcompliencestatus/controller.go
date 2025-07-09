package taxcompliancestatus

import (
	"fmt"
	"front-office/common/constant"

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
	TaxComplianceStatus(c *fiber.Ctx) error
}

func (ctrl *controller) TaxComplianceStatus(c *fiber.Ctx) error {
	reqBody := c.Locals(constant.Request).(*taxComplianceStatusRequest)
	apiKey := fmt.Sprintf("%v", c.Locals(constant.APIKey))
	memberId := fmt.Sprintf("%v", c.Locals(constant.UserId))
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	result, err := ctrl.svc.CallTaxCompliance(apiKey, memberId, companyId, reqBody)
	if err != nil {
		return err
	}

	return c.Status(result.StatusCode).JSON(result)
}
