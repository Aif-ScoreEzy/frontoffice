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
	req := c.Locals("request").(*taxComplianceStatusRequest)
	apiKey := fmt.Sprintf("%v", c.Locals("apiKey"))
	memberIdStr := fmt.Sprintf("%v", c.Locals("userId"))
	companyIdStr := fmt.Sprintf("%v", c.Locals("companyId"))

	result, err := ctrl.svc.CallTaxCompliance(apiKey, memberIdStr, companyIdStr, req)
	if err != nil {
		return err
	}

	return c.Status(result.StatusCode).JSON(result)
}
