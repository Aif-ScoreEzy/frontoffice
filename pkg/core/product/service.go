package product

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
	GetProductBySlug(slug string) (*model.AifcoreAPIResponse[productResponseData], error)
}

func (svc *service) GetProductBySlug(slug string) (*model.AifcoreAPIResponse[productResponseData], error) {
	response, err := svc.repo.CallGetProductBySlug(slug)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[productResponseData](response)
}
