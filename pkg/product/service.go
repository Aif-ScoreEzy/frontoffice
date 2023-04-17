package product

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func CreateProductSvc(req ProductRequest) (Product, error) {
	productID := uuid.NewString()
	dataProduct := Product{
		ID:      productID,
		Name:    req.Name,
		Slug:    slug.Make(req.Name),
		Version: req.Version,
		Url:     req.Url,
		Key:     req.Key,
	}

	product, err := Create(dataProduct)
	if err != nil {
		return product, err
	}

	return product, nil
}
