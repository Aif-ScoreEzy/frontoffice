package company

import (
	"fmt"
	"front-office/config/database"

	"gorm.io/gorm"
)

func FindOneByID(company Company) (Company, error) {
	err := database.DBConn.First(&company).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return company, fmt.Errorf("Company with ID %s not found", company.ID)
		}

		return company, fmt.Errorf("Failed to find company with ID %s: %v", company.ID, err)
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
