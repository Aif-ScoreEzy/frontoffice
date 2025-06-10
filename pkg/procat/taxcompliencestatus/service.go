package taxcompliancestatus

import (
	"front-office/common/model"
	"front-office/helper"
)

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

type service struct {
	repo Repository
}

type Service interface {
	CallTaxCompliance(apiKey string, request *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceDataResponse], error)
}

func (svc *service) CallTaxCompliance(apiKey string, request *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceDataResponse], error) {
	response, err := svc.repo.CallTaxComplianceStatusAPI(apiKey, request)
	if err != nil {
		return nil, err
	}

	return helper.ParseProCatAPIResponse[taxComplianceDataResponse](response)
}
