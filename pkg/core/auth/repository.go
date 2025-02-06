package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/company"
	"front-office/pkg/core/member"

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
	CreateAdmin(company *company.MstCompany, user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error)
	CreateMember(user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error)
	PasswordReset(id, token string, req *PasswordResetRequest) (*http.Response, error)
	VerifyMemberAif(req *PasswordResetRequest, memberId uint) (*http.Response, error)
	ChangePasswordAifCore(memberId string, req *ChangePasswordRequest) (*http.Response, error)
	AuthMemberAifCore(req *UserLoginRequest) (*http.Response, error)
}

func (repo *repository) CreateAdmin(company *company.MstCompany, user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error) {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		user.CompanyId = company.CompanyId
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

func (repo *repository) CreateMember(user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error) {
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

func (repo *repository) VerifyMemberAif(req *PasswordResetRequest, memberId uint) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/activation-tokens`, repo.Cfg.Env.AifcoreHost, memberId)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) ChangePasswordAifCore(memberId string, req *ChangePasswordRequest) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/change-password`, repo.Cfg.Env.AifcoreHost, memberId)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) AuthMemberAifCore(req *UserLoginRequest) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/middleware/auth-member-login"

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
