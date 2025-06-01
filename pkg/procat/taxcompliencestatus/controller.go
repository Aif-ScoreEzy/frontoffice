package taxcompliancestatus

import (
	"front-office/common/constant"
	"front-office/helper"
	"log"

	"github.com/gofiber/fiber/v2"
)

func NewController(svc Service) Controller {
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
	apiKey, _ := c.Locals("apiKey").(string)

	res, err := ctrl.svc.CallTaxCompliance(apiKey, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	if res.StatusCode >= 400 {
		log.Printf("upstream error: status=%d", res.StatusCode)
		statusCode, resp := helper.GetError(constant.UpstreamError)

		return c.Status(statusCode).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
