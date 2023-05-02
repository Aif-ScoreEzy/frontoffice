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

func GetAllProducts(c *fiber.Ctx) error {
	products, err := GetAllProductsSvc()
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Success to get all products",
		products,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := IsProductIDExistSvc(id)
	if product.Name == "" {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(resp)
	} else if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a role by ID",
		product,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateProductByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateProductRequest)
	id := c.Params("id")

	_, err := IsProductIDExistSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	product, err := UpdateProductByIDSvc(*req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResp := ProductResponse{
		Name:    product.Name,
		Slug:    product.Slug,
		Version: product.Version,
		Url:     product.Url,
		Key:     product.Key,
	}

	resp := helper.ResponseSuccess(
		"Success to update product",
		dataResp,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
