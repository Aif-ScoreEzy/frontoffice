package user

import (
	"front-office/config/database"
	"front-office/pkg/company"

	"gorm.io/gorm"
)

func FindOneByEmail(email string) (*User, error) {
	var user *User
	tx := database.DBConn.Debug().First(&user, "email = ?", email)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func FindOneByUsername(username string) (*User, error) {
	var user *User
	tx := database.DBConn.Debug().Preload("Role").First(&user, "username = ?", username)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func FindOneByKey(key string) (*User, error) {
	var user *User
	tx := database.DBConn.Debug().Preload("Role").First(&user, "key = ?", key)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func FindOneByID(id string) (*User, error) {
	var user *User
	err := database.DBConn.Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Create(company *company.Company, user *User) (*User, error) {
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

func UpdateOneByID(req *User, id string) (*User, error) {
	var user *User
	database.DBConn.Debug().First(&user, "id = ?", id)

	err := database.DBConn.Debug().Model(&user).
		Where("id = ?", id).Updates(req).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateOneByKey(key string) (*User, error) {
	var user *User
	database.DBConn.Debug().First(&user, "key = ?", key)

	err := database.DBConn.Debug().Model(&user).Update("active", true).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func DeactiveOneByEmail(email string) (*User, error) {
	var user *User
	database.DBConn.Debug().First(&user, "email = ?", email)

	err := database.DBConn.Debug().Model(&user).Update("active", false).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func FindAll() ([]*User, error) {
	var users []*User

	result := database.DBConn.Debug().Preload("Role").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
