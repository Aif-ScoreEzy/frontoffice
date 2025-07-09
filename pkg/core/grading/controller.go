package grading

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{Svc: service}
}

type controller struct {
	Svc Service
}

type Controller interface {
	CreateGradings(c *fiber.Ctx) error
	GetGradings(c *fiber.Ctx) error
	ReplaceGradings(c *fiber.Ctx) error
	ReplaceGradingsNew(c *fiber.Ctx) error
}

func (ctrl *controller) CreateGradings(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*CreateGradingsRequest)
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	var gradings []*Grading
	for _, createGradingRequest := range req.CreateGradingsRequest {
		grading, _ := ctrl.Svc.GetGradingByGradinglabelSvc(createGradingRequest.GradingLabel, companyId)
		if grading != nil {
			statusCode, res := helper.GetError(constant.DuplicateGrading)
			return c.Status(statusCode).JSON(res)
		}

		grading, err := ctrl.Svc.CreateGradingSvc(createGradingRequest, companyId)
		if err != nil {
			statusCode, res := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(res)
		}

		gradings = append(gradings, grading)
	}

	res := helper.ResponseSuccess(
		"succeed to create gradings",
		gradings,
	)

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ctrl *controller) GetGradings(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	result, err := ctrl.Svc.GetGradings(companyId)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	res := helper.ResponseSuccess(
		"succeed to get gradings",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) ReplaceGradings(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*CreateGradingsRequest)
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	err := ctrl.Svc.ReplaceAllGradingsSvc(req, companyId)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	gradings, err := ctrl.Svc.GetGradings(companyId)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	res := helper.ResponseSuccess(
		"succeed to update gradings by id",
		gradings,
	)

	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) ReplaceGradingsNew(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*CreateGradingsNewRequest)
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))

	err := ctrl.Svc.ReplaceAllGradingsNewSvc(req, companyId)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	result, err := ctrl.Svc.GetGradings(companyId)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	res := helper.ResponseSuccess(
		"succeed to update gradings by id",
		result,
	)

	return c.Status(fiber.StatusOK).JSON(res)
}
