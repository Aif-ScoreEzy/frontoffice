package taxscore

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
	CallTaxScore(apiKey string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreDataResponse], error)
}

func (svc *service) CallTaxScore(apiKey string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreDataResponse], error) {
	response, err := svc.repo.CallTaxScoreAPI(apiKey, request)
	if err != nil {
		return nil, err
	}

	return helper.ParseProCatAPIResponse[taxScoreDataResponse](response)
}
