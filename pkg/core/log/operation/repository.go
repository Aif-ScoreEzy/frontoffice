package operation

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
	FetchLogOperations(filter *LogOperationFilter) (*LogOperation, error)
	FetchByRange(filter *LogRangeFilter) (*http.Response, error)
	AddLogOperation(req *AddLogRequest) error
}

func (repo *repository) FetchLogOperations(filter *LogOperationFilter) (*LogOperation, error) {
	url := fmt.Sprintf("%s/api/middleware/auth-member-login", repo.cfg.Env.AifcoreHost)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("name", filter.Name)
	q.Add("role", filter.Role)
	q.Add("event", filter.Event)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*LogOperation](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) FetchByRange(filter *LogRangeFilter) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/operation/range"

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("company_id", filter.CompanyId)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) AddLogOperation(reqBody *AddLogRequest) error {
	url := fmt.Sprintf("%s/api/core/logging/operation", repo.cfg.Env.AifcoreHost)

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
