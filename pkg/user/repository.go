package user

import (
	"fmt"
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

func (user User) FindOneByKey() User {
	database.DBConn.Preload("Role").First(&user, "key = ?", user.Key)

	return user
}

func (user User) FindOneByID() (User, error) {
	err := database.DBConn.Preload("Role").First(&user, "id = ?", user.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, fmt.Errorf("User with ID %s not found", user.ID)
		}

		return user, fmt.Errorf("Failed to find user with ID %s: %v", user.ID, err)
	}

	return user, nil
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

func UpdateOneByID(req User, id string) (User, error) {
	var user User
	database.DBConn.Debug().First(&user, "id = ?", id)

	result := database.DBConn.Debug().Model(&user).
		Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func UpdateOneByKey(req User, key string) (User, error) {
	var user User
	database.DBConn.Debug().First(&user, "key = ?", key)

	err := database.DBConn.Debug().Model(&user).Updates(req).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func DeactiveOneByEmail(email string) (User, error) {
	var user User
	database.DBConn.Debug().First(&user, "email = ?", email)

	err := database.DBConn.Debug().Model(&user).Update("active", false).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
