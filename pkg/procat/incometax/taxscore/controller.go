package taxscore

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
	TaxScore(c *fiber.Ctx) error
}

func (ctrl *controller) TaxScore(c *fiber.Ctx) error {
	req := c.Locals("request").(*taxScoreRequest)
	apiKey, _ := c.Locals("apiKey").(string)
	memberId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	result, err := ctrl.svc.CallTaxScore(apiKey, memberId, companyId, req)
	if err != nil {
		return err
	}

	return c.Status(result.StatusCode).JSON(result)
}
