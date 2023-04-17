package product

import (
	"front-office/helper"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	req := c.Locals("request").(*ProductRequest)

	product, err := CreateProductSvc(*req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := ProductResponse{
		Name:    product.Name,
		Slug:    product.Slug,
		Version: product.Version,
		Url:     product.Url,
		Key:     product.Key,
	}

	resp := helper.ResponseSuccess(
		"Success to create product",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
