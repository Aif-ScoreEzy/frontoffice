package user

import (
	"front-office/config/database"
	"front-office/pkg/company"
	"strings"
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

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "id = ?", id).Error
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
		return nil, errTx
	}

	database.DBConn.Preload("Company").Preload("Company.Industry").Preload("Role").Preload("Role.Permissions").First(&user)

	return user, nil
}

func CreateMember(user *User, activationToken *ActivationToken) (*User, error) {
	errTx := database.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Debug().Create(&activationToken).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return nil, errTx
	}

	database.DBConn.Preload("Company").Preload("Role").First(&user)

	return user, nil
}

func CreateActivationToken(data *ActivationToken) error {
	err := database.DBConn.Debug().Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func VerifyUserTx(req map[string]interface{}, id, token string) (*User, error) {
	errTX := database.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Model(&ActivationToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		if err := tx.Debug().Model(&User{}).Where("id = ?", id).Updates(req).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return nil, errTX
	}

	userDetail, err := FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return userDetail, nil
}

func UpdateOneByID(req map[string]interface{}, id, companyID string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Model(&user).
		Where("id = ? AND company_id = ?", id, companyID).Updates(req).Error
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

func FindAll(limit, offset int, keyword, roleID, active, startTime, endTime, companyID string) ([]User, error) {
	var users []User

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := database.DBConn.Debug().Preload("Role").Where("company_id = ? AND LOWER(name) LIKE ?", companyID, "%"+keywordToLower+"%")

	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	if active != "" {
		query = query.Where("active = ?", active)
	}

	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

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

func GetTotalData(keyword, roleID, active, startTime, endTime, companyID string) (int64, error) {
	var users []User
	var count int64

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := database.DBConn.Debug().Where("company_id = ? AND LOWER(name) LIKE ?", companyID, "%"+keywordToLower+"%")
	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	if active != "" {
		query = query.Where("active = ?", active)
	}

	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	err := query.Find(&users).Count(&count).Error

	return count, err
}
