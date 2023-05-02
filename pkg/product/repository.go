package product

import (
	"fmt"
	"front-office/config/database"

	"gorm.io/gorm"
)

func Create(product Product) (Product, error) {
	err := database.DBConn.Debug().Create(&product).Error
	if err != nil {
		return product, err
	}

	database.DBConn.Debug().First(&product)

	return product, nil
}

func FindAll() ([]Product, error) {
	var products []Product

	result := database.DBConn.Debug().Find(&products)
	if result.Error != nil {
		return products, result.Error
	}

	return products, nil
}

func FindOneByID(product Product) (Product, error) {
	err := database.DBConn.Debug().First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return product, fmt.Errorf("Product with ID %s not found", product.ID)
		}

		return product, fmt.Errorf("Failed to find role with ID %s: %v", product.ID, err)
	}

	return product, nil
}

func UpdateOneByID(req Product, id string) (Product, error) {
	var product Product
	database.DBConn.Debug().First(&product, "id = ?", id)

	err := database.DBConn.Debug().Model(&product).Where("id = ?", id).Updates(req).Error
	if err != nil {
		return product, err
	}

	return product, nil
}
