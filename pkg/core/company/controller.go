package company

import (
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/industry"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service, svcIndustry industry.Service) Controller {
	return &controller{Svc: service, SvcIndustry: svcIndustry}
}

type controller struct {
	Svc         Service
	SvcIndustry industry.Service
}

type Controller interface {
	UpdateCompanyById(c *fiber.Ctx) error
}

func (ctrl *controller) UpdateCompanyById(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*UpdateCompanyRequest)
	id := c.Params("id")

	if req.IndustryId != "" {
		_, err := ctrl.SvcIndustry.IsIndustryIdExistSvc(req.IndustryId)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}
	}

	company, err := ctrl.Svc.UpdateCompanyByIdSvc(*req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UpdateCompanyResponse{
		Id:              company.Id,
		CompanyName:     company.CompanyName,
		CompanyAddress:  company.CompanyAddress,
		CompanyPhone:    company.CompanyPhone,
		AgreementNumber: company.AgreementNumber,
		PaymentScheme:   company.PaymentScheme,
		PostpaidActive:  company.PostpaidActive,
		IndustryId:      company.IndustryId,
	}

	resp := helper.ResponseSuccess(
		"Success to update company",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
