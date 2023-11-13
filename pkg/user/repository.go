package user

import (
	"front-office/config/database"
	"strings"
	"time"
)

func FindOneByEmail(email string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindOneByUserID(id string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "id = ?", id).Error
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

func FindOneByUserIDAndCompanyID(id, companyID string) (*User, error) {
	var user *User

	err := database.DBConn.Debug().Preload("Role").Preload("Company").First(&user, "id = ? AND company_id = ?", id, companyID).Error
	if err != nil {
		return nil, err
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
		query = query.Where("status = ?", status)
	}

	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	err := query.Find(&users).Count(&count).Error

	return count, err
}
