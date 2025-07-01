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
	CallVerifyMemberAPI(userId uint, req *PasswordResetRequest) error
	CallChangePasswordAPI(userId string, req *ChangePasswordRequest) error
	CallPasswordResetAPI(userId uint, token string, req *PasswordResetRequest) error
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

func (repo *repository) CallVerifyMemberAPI(userId uint, reqBody *PasswordResetRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/activation-tokens`, repo.cfg.Env.AifcoreHost, userId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallPasswordResetAPI(userId uint, token string, reqBody *PasswordResetRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/password-reset-tokens/%v`, repo.cfg.Env.AifcoreHost, userId, token)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallChangePasswordAPI(userId string, reqBody *ChangePasswordRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/change-password`, repo.cfg.Env.AifcoreHost, userId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) AuthMemberAifCore(reqBody *userLoginRequest) (*loginResponseData, error) {
	url := fmt.Sprintf("%s/api/middleware/auth-member-login", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	// send http request
	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	// parse structured response
	apiResp, err := helper.ParseAifcoreAPIResponse[*loginResponseData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}
