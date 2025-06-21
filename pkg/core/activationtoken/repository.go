package activationtoken

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
	CallGetActivationTokenAPI(token string) (*MstActivationToken, error)
	CallCreateActivationTokenAPI(memberId string, req *CreateActivationTokenRequest) error
}

func (repo *repository) CallGetActivationTokenAPI(token string) (*MstActivationToken, error) {
	url := fmt.Sprintf(`%v/api/core/member/activation-tokens/%v`, repo.cfg.Env.AifcoreHost, token)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*MstActivationToken](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) CallCreateActivationTokenAPI(memberId string, reqBody *CreateActivationTokenRequest) error {
	url := fmt.Sprintf(`%v/api/core/member/%v/activation-tokens`, repo.cfg.Env.AifcoreHost, memberId)

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
