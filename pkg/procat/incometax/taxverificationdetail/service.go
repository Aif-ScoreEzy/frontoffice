package taxverificationdetail

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
	CallTaxVerification(apiKey string, request *taxVerificationRequest) (*model.ProCatAPIResponse[taxVerificationDataResponse], error)
}

func (svc *service) CallTaxVerification(apiKey string, request *taxVerificationRequest) (*model.ProCatAPIResponse[taxVerificationDataResponse], error) {
	response, err := svc.repo.CallTaxVerificationAPI(apiKey, request)
	if err != nil {
		return nil, err
	}

	return helper.ParseProCatAPIResponse[taxVerificationDataResponse](response)
}
