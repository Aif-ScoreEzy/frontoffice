package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{DB: db, Cfg: cfg}
}

type repository struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	FindOneByEmail(email string) (*User, error)
	FindOneByUserId(id string) (*User, error)
	FindOneByKey(key string) (*User, error)
	FindOneByUserIdAndCompanyId(id, companyId string) (*User, error)
	UpdateOneById(req map[string]interface{}, user *User) (*User, error)
	FindAll(limit, offset int, keyword, roleId, status, startTime, endTime, companyId string) ([]User, error)
	DeleteById(id string) error
	GetTotalData(keyword, roleId, status, startTime, endTime, companyId string) (int64, error)
	AddMemberAifCore(req *RegisterMemberRequest) (*http.Response, error)
	FindOneAifCore(query *FindUserQuery) (*http.Response, error)
	UpdateOneByIdAifCore(req map[string]interface{}, memberId uint) (*http.Response, error)
}

func (repo *repository) FindOneByEmail(email string) (*User, error) {
	var user *User

	err := repo.DB.Preload("Role").Preload("Company").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOneByUserId(id string) (*User, error) {
	var user *User

	err := repo.DB.Preload("Role").Preload("Company").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOneByKey(key string) (*User, error) {
	var user *User

	err := repo.DB.Preload("Role").Preload("Company").First(&user, "key = ?", key).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOneByUserIdAndCompanyId(id, companyId string) (*User, error) {
	var user *User

	err := repo.DB.Preload("Role").Preload("Company").First(&user, "id = ? AND company_id = ?", id, companyId).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) UpdateOneById(req map[string]interface{}, user *User) (*User, error) {
	err := repo.DB.Model(&user).
		Where("id = ? AND company_id = ?", user.Id, user.CompanyId).Updates(req).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindAll(limit, offset int, keyword, roleId, status, startTime, endTime, companyId string) ([]User, error) {
	var users []User

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := repo.DB.Preload("Role").Where("company_id = ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", companyId, "%"+keywordToLower+"%", "%"+keywordToLower+"%")

	if roleId != "" {
		query = query.Where("role_id = ?", roleId)
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

func (repo *repository) DeleteById(id string) error {
	err := repo.DB.Model(&User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetTotalData(keyword, roleId, status, startTime, endTime, companyId string) (int64, error) {
	var users []User
	var count int64

	// avoid case sensitive (uppercase/lowercase) keywords
	keywordToLower := strings.ToLower(keyword)

	query := repo.DB.Where("company_id = ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", companyId, "%"+keywordToLower+"%", "%"+keywordToLower+"%")
	if roleId != "" {
		query = query.Where("role_id = ?", roleId)
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

func (repo *repository) AddMemberAifCore(req *RegisterMemberRequest) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/member/addmember"

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	request.Header.Set(constant.XAPIKey, repo.Cfg.Env.XModuleKey)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) FindOneAifCore(query *FindUserQuery) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/by`, repo.Cfg.Env.AifcoreHost)

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("id", query.Id)
	q.Add("email", query.Email)
	q.Add("username", query.Username)
	q.Add("key", query.Key)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) UpdateOneByIdAifCore(req map[string]interface{}, memberId uint) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%d`, repo.Cfg.Env.AifcoreHost, memberId)

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
