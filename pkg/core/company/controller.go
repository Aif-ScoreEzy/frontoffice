package company

import (
	"front-office/helper"
	"front-office/pkg/core/industry"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{Svc: service}
}

type controller struct {
	Svc         Service
	SvcIndustry industry.Service
}

type Controller interface {
	UpdateCompanyByID(c *fiber.Ctx) error
}

func (ctrl *controller) UpdateCompanyByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateCompanyRequest)
	id := c.Params("id")

	_, err := ctrl.SvcIndustry.IsIndustryIDExistSvc(req.IndustryID)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	company, err := ctrl.Svc.UpdateCompanyByIDSvc(*req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UpdateCompanyResponse{
		ID:              company.ID,
		CompanyName:     company.CompanyName,
		CompanyAddress:  company.CompanyAddress,
		CompanyPhone:    company.CompanyPhone,
		AgreementNumber: company.AgreementNumber,
		PaymentScheme:   company.PaymentScheme,
		PostpaidActive:  company.PostpaidActive,
		IndustryID:      company.IndustryID,
	}

	resp := helper.ResponseSuccess(
		"Success to update company",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
