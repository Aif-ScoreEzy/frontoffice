package template

import (
	"front-office/common/constant"
	"front-office/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	ListTemplates(c *fiber.Ctx) error
	DownloadTemplate(c *fiber.Ctx) error
}

type controller struct {
	Svc Service
}

func NewController(service Service) Controller {
	return &controller{Svc: service}
}
func (ctrl *controller) ListTemplates(c *fiber.Ctx) error {
	templates, err := ctrl.Svc.ListTemplates()
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(resp)
	}

	return c.JSON(templates)
}

// Download specific template
func (ctrl *controller) DownloadTemplate(c *fiber.Ctx) error {
	var req DownloadRequest
	if err := c.QueryParser(&req); err != nil {
		_, resp := helper.GetError(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if req.Category == "" {
		_, resp := helper.GetError("category parameter is required")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	if req.Filename == "" {
		req.Filename = "template.csv"
	} else if !strings.HasSuffix(req.Filename, ".csv") {
		req.Filename += ".csv"
	}

	path, err := ctrl.Svc.DownloadTemplate(req)
	if err != nil {
		statusCode, resp := helper.GetError(constant.TemplateNotFound)

		return c.Status(statusCode).JSON(resp)
	}

	return c.Download(path)
}
