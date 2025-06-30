package passwordresettoken

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
	CallCreatePasswordResetTokenAPI(userId string, reqBody *CreatePasswordResetTokenRequest) error
	CallGetPasswordResetTokenAPI(token string) (*MstPasswordResetToken, error)
	CallDeletePasswordResetTokenAPI(id string) error
}

func (repo *repository) CallGetPasswordResetTokenAPI(token string) (*MstPasswordResetToken, error) {
	url := fmt.Sprintf(`%v/api/core/member/password-reset-tokens/%v`, repo.cfg.Env.AifcoreHost, token)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*MstPasswordResetToken](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) CallCreatePasswordResetTokenAPI(userId string, reqBody *CreatePasswordResetTokenRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/password-reset-tokens`, repo.cfg.Env.AifcoreHost, userId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallDeletePasswordResetTokenAPI(id string) error {
	url := fmt.Sprintf(`%v/api/core/member/password-reset-tokens/%v`, repo.cfg.Env.AifcoreHost, id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
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
