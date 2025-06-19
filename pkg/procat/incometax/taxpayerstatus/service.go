package taxpayerstatus

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
	CallTaxPayerStatus(apiKey string, request *taxPayerStatusRequest) (*model.ProCatAPIResponse[taxPayerStatusDataResponse], error)
}

func (svc *service) CallTaxPayerStatus(apiKey string, request *taxPayerStatusRequest) (*model.ProCatAPIResponse[taxPayerStatusDataResponse], error) {
	response, err := svc.repo.CallTaxPayerStatusAPI(apiKey, request)
	if err != nil {
		return nil, err
	}

	return helper.ParseProCatAPIResponse[taxPayerStatusDataResponse](response)
}
