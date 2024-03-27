package user

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

type repository struct {
	DB *gorm.DB
}

type Repository interface {
	FindOneByEmail(email string) (*User, error)
	FindOneByUserID(id string) (*User, error)
	FindOneByKey(key string) (*User, error)
	FindOneByUserIDAndCompanyID(id, companyID string) (*User, error)
	UpdateOneByID(req map[string]interface{}, user *User) (*User, error)
	FindAll(limit, offset int, keyword, roleID, status, startTime, endTime, companyID string) ([]User, error)
	DeleteByID(id string) error
	GetTotalData(keyword, roleID, status, startTime, endTime, companyID string) (int64, error)
}

func (repo *repository) FindOneByEmail(email string) (*User, error) {
	var user *User

	err := repo.DB.Debug().Preload("Role").Preload("Company").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOneByUserID(id string) (*User, error) {
	var user *User

	err := repo.DB.Debug().Preload("Role").Preload("Company").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOneByKey(key string) (*User, error) {
	var user *User

	err := repo.DB.Debug().Preload("Role").Preload("Company").First(&user, "key = ?", key).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOneByUserIDAndCompanyID(id, companyID string) (*User, error) {
	var user *User

	err := repo.DB.Debug().Preload("Role").Preload("Company").First(&user, "id = ? AND company_id = ?", id, companyID).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) UpdateOneByID(req map[string]interface{}, user *User) (*User, error) {
	err := repo.DB.Debug().Model(&user).
		Where("id = ? AND company_id = ?", user.ID, user.CompanyID).Updates(req).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindAll(limit, offset int, keyword, roleID, status, startTime, endTime, companyID string) ([]User, error) {
	var users []User

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := repo.DB.Debug().Preload("Role").Where("company_id = ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", companyID, "%"+keywordToLower+"%", "%"+keywordToLower+"%")

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

func (repo *repository) DeleteByID(id string) error {
	err := repo.DB.Debug().Model(&User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetTotalData(keyword, roleID, status, startTime, endTime, companyID string) (int64, error) {
	var users []User
	var count int64

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := repo.DB.Debug().Where("company_id = ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", companyID, "%"+keywordToLower+"%", "%"+keywordToLower+"%")
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
