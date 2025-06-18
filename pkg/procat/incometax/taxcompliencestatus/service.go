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
	CallTaxCompliance(apiKey, jobId string, request *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceDataResponse], error)
}

func (svc *service) CallTaxCompliance(apiKey, jobId string, request *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceDataResponse], error) {
	response, err := svc.repo.CallTaxComplianceStatusAPI(apiKey, jobId, request)
	if err != nil {
		return nil, err
	}

	return helper.ParseProCatAPIResponse[taxComplianceDataResponse](response)
}
