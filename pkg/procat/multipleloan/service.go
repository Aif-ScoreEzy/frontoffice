package multipleloan

import (
	"front-office/app/config"
	"front-office/common/model"
	"front-office/helper"
)

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{
		Cfg:  cfg,
		Repo: repo,
	}
}

type service struct {
	Cfg  *config.Config
	Repo Repository
}

type Service interface {
	CallMultipleLoan7Days(request *multipleLoanRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error)
	CallMultipleLoan30Days(request *multipleLoanRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error)
	CallMultipleLoan90Days(request *multipleLoanRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error)
}

func (svc *service) CallMultipleLoan7Days(request *multipleLoanRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error) {
	response, err := svc.Repo.CallMultipleLoan7Days(request, apiKey, memberId, companyId)
	if err != nil {
		return nil, err
	}

	result, err := helper.ParseProCatAPIResponse[dataMultipleLoanResponse](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) CallMultipleLoan30Days(request *multipleLoanRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error) {
	response, err := svc.Repo.CallMultipleLoan30Days(request, apiKey, memberId, companyId)
	if err != nil {
		return nil, err
	}

	result, err := helper.ParseProCatAPIResponse[dataMultipleLoanResponse](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) CallMultipleLoan90Days(request *multipleLoanRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error) {
	response, err := svc.Repo.CallMultipleLoan90Days(request, apiKey, memberId, companyId)
	if err != nil {
		return nil, err
	}

	result, err := helper.ParseProCatAPIResponse[dataMultipleLoanResponse](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
