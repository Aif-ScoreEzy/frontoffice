package taxscore

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
	CallTaxScore(apiKey string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreDataResponse], error)
}

func (svc *service) CallTaxScore(apiKey string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreDataResponse], error) {
	response, err := svc.repo.CallTaxScoreAPI(apiKey, request)
	if err != nil {
		return nil, err
	}

	return helper.ParseProCatAPIResponse[taxScoreDataResponse](response)
}
