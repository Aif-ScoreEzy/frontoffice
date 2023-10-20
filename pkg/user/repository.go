package user

import (
	"front-office/config/database"
	"strings"
	"time"

	"gorm.io/gorm"
)

func FindOneByEmail(email string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindOneByKey(key string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "key = ?", key).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindOneByID(id, companyID string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "id = ? AND company_id = ?", id, companyID).Error
	if err != nil {
		return nil, err
	}

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

	return user, nil
}

func FindOneActivationTokenBytoken(token string) (*ActivationToken, error) {
	var activationToken *ActivationToken

	err := database.DBConn.Debug().First(&activationToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func FindOneActivationTokenByUserID(userID string) (*ActivationToken, error) {
	var activationToken *ActivationToken

	err := database.DBConn.Debug().First(&activationToken, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func CreateActivationToken(activationToken *ActivationToken) (*ActivationToken, error) {
	err := database.DBConn.Debug().Create(&activationToken).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func VerifyUserTx(req map[string]interface{}, userID, token string) (*User, error) {
	var user *User
	errTX := database.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Model(&ActivationToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		if err := tx.Debug().Model(&user).Where("id = ?", userID).Updates(req).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return nil, errTX
	}

	return user, nil
}

func UpdateOneByID(req map[string]interface{}, user *User) (*User, error) {
	err := database.DBConn.Debug().Model(&user).
		Where("id = ? AND company_id = ?", user.ID, user.CompanyID).Updates(req).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindAll(limit, offset int, keyword, roleID, status, startTime, endTime, companyID string) ([]User, error) {
	var users []User

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := database.DBConn.Debug().Preload("Role").Where("company_id = ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", companyID, "%"+keywordToLower+"%", "%"+keywordToLower+"%")

	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
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

func GetTotalData(keyword, roleID, status, startTime, endTime, companyID string) (int64, error) {
	var users []User
	var count int64

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := database.DBConn.Debug().Where("company_id = ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", companyID, "%"+keywordToLower+"%", "%"+keywordToLower+"%")
	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	if status != "" {
		query = query.Where("active = ?", status)
	}

	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	err := query.Find(&users).Count(&count).Error

	return count, err
}
