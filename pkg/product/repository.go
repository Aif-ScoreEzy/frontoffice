package product

import "front-office/config/database"

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
