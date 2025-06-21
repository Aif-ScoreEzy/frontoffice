package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/httpclient"

	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{cfg, client}
}

type repository struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type Repository interface {
	// CreateAdmin(company *company.MstCompany, user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error)
	// CreateMember(user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error)
	CallVerifyMemberAPI(memberId uint, req *PasswordResetRequest) error
	PasswordReset(memberId uint, token string, req *PasswordResetRequest) (*http.Response, error)
	ChangePasswordAifCore(memberId string, req *ChangePasswordRequest) (*http.Response, error)
	AuthMemberAifCore(req *userLoginRequest) (*loginResponseData, error)
}

// func (repo *repository) CreateAdmin(company *company.MstCompany, user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error) {
// 	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Create(&company).Error; err != nil {
// 			return err
// 		}

// 		user.CompanyId = company.CompanyId
// 		if err := tx.Create(&user).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Create(&activationToken).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if errTx != nil {
// 		return user, errTx
// 	}

// 	repo.DB.Preload("Company").Preload("Company.Industry").Preload("Role").First(&user)

// 	return user, errTx
// }

// func (repo *repository) CreateMember(user *member.MstMember, activationToken *activationtoken.MstActivationToken) (*member.MstMember, error) {
// 	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Create(&user).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Create(&activationToken).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if errTx != nil {
// 		return nil, errTx
// 	}

// 	return user, nil
// }

func (repo *repository) CallVerifyMemberAPI(memberId uint, reqBody *PasswordResetRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/activation-tokens`, repo.cfg.Env.AifcoreHost, memberId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) PasswordReset(memberId uint, token string, req *PasswordResetRequest) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/password-reset-tokens/%v`, repo.cfg.Env.AifcoreHost, memberId, token)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) ChangePasswordAifCore(memberId string, req *ChangePasswordRequest) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/change-password`, repo.cfg.Env.AifcoreHost, memberId)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) AuthMemberAifCore(reqBody *userLoginRequest) (*loginResponseData, error) {
	url := fmt.Sprintf("%s/api/middleware/auth-member-login", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	// send http request
	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// parse structured response
	apiResp, err := helper.ParseAifcoreAPIResponse[*loginResponseData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}
