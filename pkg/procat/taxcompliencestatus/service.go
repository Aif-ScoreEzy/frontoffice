package taxcompliancestatus

import (
	"front-office/app/config"
	"front-office/common/model"
	"front-office/helper"
)

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{
		cfg,
		repo,
	}
}

type service struct {
	cfg  *config.Config
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
