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
