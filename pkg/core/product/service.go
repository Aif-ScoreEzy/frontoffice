package product

import (
	"front-office/internal/apperror"
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
	GetProductBySlug(slug string) (*productResponseData, error)
}

func (svc *service) GetProductBySlug(slug string) (*productResponseData, error) {
	product, err := svc.repo.GetProductAPI(slug)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch product")
	}
	if product.ProductId == 0 {
		return nil, apperror.NotFound("product not found")
	}

	return product, err
}
