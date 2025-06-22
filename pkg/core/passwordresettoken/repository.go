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
	CallCreatePasswordResetToken(userId string, reqBody *CreatePasswordResetTokenRequest) error
	FindOnePasswordResetTokenByToken(token string) (*http.Response, error)
	DeletePasswordResetToken(id string) (*http.Response, error)
}

func (repo *repository) FindOnePasswordResetTokenByToken(token string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/password-reset-tokens/%v`, repo.cfg.Env.AifcoreHost, token)

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) CallCreatePasswordResetToken(userId string, reqBody *CreatePasswordResetTokenRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/password-reset-tokens`, repo.cfg.Env.AifcoreHost, userId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) DeletePasswordResetToken(id string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/password-reset-tokens/%v`, repo.cfg.Env.AifcoreHost, id)

	request, _ := http.NewRequest(http.MethodDelete, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
