package taxcompliancestatus

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
	TaxComplianceStatus(c *fiber.Ctx) error
}

func (ctrl *controller) TaxComplianceStatus(c *fiber.Ctx) error {
	reqBody := c.Locals("request").(*taxComplianceStatusRequest)
	apiKey := fmt.Sprintf("%v", c.Locals("apiKey"))
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	result, err := ctrl.svc.CallTaxCompliance(apiKey, memberId, companyId, reqBody)
	if err != nil {
		return err
	}

	return c.Status(result.StatusCode).JSON(result)
}
