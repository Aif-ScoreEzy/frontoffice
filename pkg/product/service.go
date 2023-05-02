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

func GetAllProductsSvc() ([]Product, error) {
	products, err := FindAll()
	if err != nil {
		return products, err
	}

	return products, nil
}

func IsProductIDExistSvc(id string) (Product, error) {
	product := Product{
		ID: id,
	}

	result, err := FindOneByID(product)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateProductByIDSvc(req UpdateProductRequest, id string) (Product, error) {
	dataReq := Product{
		Name:    req.Name,
		Slug:    slug.Make(req.Name),
		Version: req.Version,
		Url:     req.Url,
	}

	product, err := UpdateOneByID(dataReq, id)
	if err != nil {
		return product, err
	}

	return product, nil
}
