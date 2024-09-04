package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/company"
	"front-office/pkg/core/user"
	"net/http"

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
	CreateAdmin(company *company.Company, user *user.User, activationToken *activationtoken.MstActivationToken) (*user.User, error)
	CreateMember(user *user.User, activationToken *activationtoken.MstActivationToken) (*user.User, error)
	PasswordReset(id, token string, req *PasswordResetRequest) (*http.Response, error)
	VerifyUserTx(req map[string]interface{}, userId, token string) (*user.User, error)
	LoginAifCoreService(req *UserLoginRequest) (*http.Response, error)
	ChangePasswordAifCoreService(req *ChangePasswordRequest) (*http.Response, error)
}

func (repo *repository) CreateAdmin(company *company.Company, user *user.User, activationToken *activationtoken.MstActivationToken) (*user.User, error) {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		user.CompanyId = company.Id
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&activationToken).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return user, errTx
	}

	repo.DB.Preload("Company").Preload("Company.Industry").Preload("Role").First(&user)

	return user, errTx
}

func (repo *repository) CreateMember(user *user.User, activationToken *activationtoken.MstActivationToken) (*user.User, error) {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&activationToken).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return nil, errTx
	}

	return user, nil
}

func (repo *repository) PasswordReset(memberId, token string, req *PasswordResetRequest) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/password-reset-tokens/%v`, repo.Cfg.Env.AifcoreHost, memberId, token)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) VerifyUserTx(req map[string]interface{}, userId, token string) (*user.User, error) {
	var user *user.User

	errTX := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&activationtoken.MstActivationToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Where("id = ?", userId).Updates(req).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return nil, errTX
	}

	return user, nil
}

func (repo *repository) LoginAifCoreService(req *UserLoginRequest) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/member/login"

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) ChangePasswordAifCoreService(req *ChangePasswordRequest) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/member/change-password"

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
