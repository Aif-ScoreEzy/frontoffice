package multipleloan

import (
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"net/http"
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repository
}

type Service interface {
	MultipleLoan(apiKey, jobId, productSlug, memberId, companyId string, request *multipleLoanRequest) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error)
}

func (svc *service) MultipleLoan(apiKey, jobId, productSlug, memberId, companyId string, request *multipleLoanRequest) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error) {
	var response *http.Response
	var err error

	switch productSlug {
	case constant.SlugMultipleLoan7Days:
		response, err = svc.repo.CallMultipleLoan7Days(request, apiKey, memberId, jobId, companyId)
		if err != nil {
			return nil, err
		}
	case constant.SlugMultipleLoan30Days:
		response, err = svc.repo.CallMultipleLoan30Days(request, apiKey, memberId, jobId, companyId)
		if err != nil {
			return nil, err
		}
	case constant.SlugMultipleLoan90Days:
		response, err = svc.repo.CallMultipleLoan90Days(request, apiKey, memberId, jobId, companyId)
		if err != nil {
			return nil, err
		}
	}

	result, err := helper.ParseProCatAPIResponse[dataMultipleLoanResponse](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
