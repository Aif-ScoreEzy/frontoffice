package grading

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func CreateGradings(c *fiber.Ctx) error {
	req := c.Locals("request").(*CreateGradingsRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	var gradings []*Grading
	for _, createGradingRequest := range req.CreateGradingsRequest {
		grading, _ := GetGradingByGradinglabelSvc(createGradingRequest.GradingLabel, companyID)
		if grading != nil {
			statusCode, res := helper.GetError(constant.DuplicateGrading)
			return c.Status(statusCode).JSON(res)
		}

		grading, err := CreateGradingSvc(createGradingRequest, companyID)
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

func GetGradings(c *fiber.Ctx) error {
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	gradings, err := GetGradingsSvc(companyID)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	res := helper.ResponseSuccess(
		"succeed to get gradings",
		gradings,
	)

	return c.Status(fiber.StatusOK).JSON(res)
}

func UpdateGradingsByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*CreateGradingsRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	err := DeleteGradingsSvc(companyID)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(res)
	}

	for _, createGradingRequest := range req.CreateGradingsRequest {
		grading, _ := GetGradingByGradinglabelSvc(createGradingRequest.GradingLabel, companyID)
		if grading != nil {
			statusCode, res := helper.GetError(constant.DuplicateGrading)
			return c.Status(statusCode).JSON(res)
		}

		_, err := CreateGradingSvc(createGradingRequest, companyID)
		if err != nil {
			statusCode, res := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(res)
		}
	}

	gradings, err := GetGradingsSvc(companyID)
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