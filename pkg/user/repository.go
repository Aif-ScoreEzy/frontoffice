package user

import (
	"front-office/config/database"
	"front-office/pkg/company"

	"gorm.io/gorm"
)

func (user User) FindOneByEmail() User {
	database.DBConn.First(&user, "email = ?", user.Email)

	return user
}

func (user User) FindOneByUsername() User {
	database.DBConn.Preload("Role").First(&user, "username = ?", user.Username)

	return user
}

func Create(company company.Company, user User) (User, error) {
	errTx := database.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		user.CompanyID = company.ID

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return user, errTx
	}

	database.DBConn.Preload("Company").Preload("Company.Industry").Preload("Role").Preload("Role.Permissions").First(&user)

	return user, errTx
}
