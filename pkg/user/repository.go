package user

import (
	"front-office/config/database"
	"front-office/pkg/company"
	"time"

	"gorm.io/gorm"
)

func FindOneByEmail(email string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").Preload("Company.Industry").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindOneByKey(key string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").Preload("Company.Industry").First(&user, "key = ?", key).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindOneByID(id string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").Preload("Company.Industry").First(&user, "id = ?", id).Error
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

func CreateMember(user *User) (*User, error) {
	err := database.DBConn.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}

	database.DBConn.Debug().Preload("Company").Preload("Role").Preload("Role.Permissions").First(&user)

	return user, nil
}

func UpdateOneByID(req *User, id string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Model(&user).
		Where("id = ?", id).Updates(req).Error
	if err != nil {
		return nil, err
	}

	userDetail, err := FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return userDetail, nil
}

func UpdateOneByKey(key string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Model(&user).Where("key = ?", key).Update("active", true).Error
	if err != nil {
		return user, err
	}

	if err := database.DBConn.First(&user, "key = ?", key).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func DeactiveOneByEmail(email string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Model(&user).Where("email = ?", email).Update("active", false).Error
	if err != nil {
		return user, err
	}

	if err := database.DBConn.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func FindAll(limit, offset int, keyword, companyID string) ([]User, error) {
	var users []User

	// result := database.DBConn.Debug().Preload("Role").Limit(limit).Offset(offset).Find(&users, "company_id = ?", companyID)

	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	// return users, nil

	// new queries
	conditions := map[string]interface{}{
		"name":  "%" + keyword + "%", // Example name condition
		"email": "%" + keyword + "%", // Example email condition
	}
	query := database.DBConn.Debug().Preload("Role")

	for key, value := range conditions {
		query = query.Or(key+" LIKE ?", value)
	}
	query = query.Where("company_id = ?", companyID)

	result := query.Limit(limit).Offset(offset).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func DeleteByID(id string) error {
	err := database.DBConn.Debug().Model(&User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func GetTotalData(keyword, companyID string) (int64, error) {
	var users []User
	var count int64

	// make sure main query is the same as FindAll method
	conditions := map[string]interface{}{
		"name":  "%" + keyword + "%", // Example name condition
		"email": "%" + keyword + "%", // Example email condition
	}

	query := database.DBConn.Debug().Preload("Role")
	for key, value := range conditions {
		query = query.Or(key+" LIKE ?", value)
	}
	query = query.Where("company_id = ?", companyID)

	err := query.Find(&users).Count(&count).Error

	return count, err
}
