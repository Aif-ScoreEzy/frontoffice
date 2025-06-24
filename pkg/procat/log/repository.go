package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/httpclient"
	"net/http"
	"time"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{
		cfg:    cfg,
		client: client,
	}
}

type repository struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type Repository interface {
	CallCreateProCatJobAPI(reqBody *CreateJobRequest) (*createJobDataResponse, error)
	CallUpdateJobAPI(jobId string, req map[string]interface{}) error
	CallGetProCatJobAPI(filter *logFilter) (*http.Response, error)
	CallGetProCatJobDetailAPI(filter *logFilter) (*http.Response, error)
}

func (repo *repository) CallCreateProCatJobAPI(reqBody *CreateJobRequest) (*createJobDataResponse, error) {
	url := fmt.Sprintf("%s/api/core/product/jobs", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set("X-Member-ID", reqBody.MemberId)
	req.Header.Set("X-Company-ID", reqBody.CompanyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*createJobDataResponse](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) CallUpdateJobAPI(jobId string, reqBody map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/core/product/jobs/%s", repo.cfg.Env.AifcoreHost, jobId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*createJobDataResponse](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallGetProCatJobAPI(filter *logFilter) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/product/" + filter.ProductSlug + "/jobs"

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	q := httpRequest.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	httpRequest.URL.RawQuery = q.Encode()

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallGetProCatJobDetailAPI(filter *logFilter) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/product/" + filter.ProductSlug + "/jobs/" + filter.JobId

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	client := http.Client{}

	return client.Do(httpRequest)
}
