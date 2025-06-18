package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
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
	CallCreateProCatJobAPI(req *CreateJobRequest) (*http.Response, error)
	CallUpdateJobAPI(jobId string, req map[string]interface{}) (*http.Response, error)
	CallGetProCatJobAPI(filter *logFilter) (*http.Response, error)
	CallGetProCatJobDetailAPI(filter *logFilter) (*http.Response, error)
}

func (repo *repository) CallCreateProCatJobAPI(req *CreateJobRequest) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/product/jobs"

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", req.MemberId)
	httpRequest.Header.Set("X-Company-ID", req.CompanyId)

	resp, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

func (repo *repository) CallUpdateJobAPI(jobId string, req map[string]interface{}) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/product/jobs/" + jobId

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPut, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
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
