package multipleloan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/internal/httpclient"
	"net/http"
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
	CallMultipleLoan7Days(request *multipleLoanRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error)
	CallMultipleLoan30Days(request *multipleLoanRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error)
	CallMultipleLoan90Days(request *multipleLoanRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error)
}

func (repo *repository) CallMultipleLoan7Days(request *multipleLoanRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/7-days"

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-API-KEY", apiKey)
	httpRequest.Header.Set("X-Member-ID", memberId)
	httpRequest.Header.Set("X-Company-ID", companyId)

	q := httpRequest.URL.Query()
	q.Add("job_id", jobId)
	httpRequest.URL.RawQuery = q.Encode()

	response, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}

func (repo *repository) CallMultipleLoan30Days(request *multipleLoanRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/30-days"

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-API-Key", apiKey)
	httpRequest.Header.Set("X-Member-ID", memberId)
	httpRequest.Header.Set("X-Company-ID", companyId)

	q := httpRequest.URL.Query()
	q.Add("job_id", jobId)
	httpRequest.URL.RawQuery = q.Encode()

	response, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}

func (repo *repository) CallMultipleLoan90Days(request *multipleLoanRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/90-days"

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-API-Key", apiKey)
	httpRequest.Header.Set("X-Member-ID", memberId)
	httpRequest.Header.Set("X-Company-ID", companyId)

	q := httpRequest.URL.Query()
	q.Add("job_id", jobId)
	httpRequest.URL.RawQuery = q.Encode()

	response, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}
