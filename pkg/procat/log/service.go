package log

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
	GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
}

func (svc *service) GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	response, err := svc.Repo.CallProCatJobAPI(filter)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}

func (svc *service) GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	response, err := svc.Repo.CallGetProCatJobDetailAPI(filter)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}
