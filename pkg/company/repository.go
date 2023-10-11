package company

import (
	"front-office/config/database"
)

func FindOneByID(id string) (*Company, error) {
	var company *Company

	err := database.DBConn.First(&company, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return company, nil
}

func UpdateOneByID(req Company, id string) (Company, error) {
	var company Company
	database.DBConn.First(&company, "id = ?", id)

	err := database.DBConn.Debug().Model(&company).Updates(req).Error
	if err != nil {
		return company, err
	}

	return company, nil
}
