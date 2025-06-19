package loanrecordchecker

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
	CallLoanRecordCheckerAPI(request *LoanRecordCheckerRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error)
}

func (repo *repository) CallLoanRecordCheckerAPI(request *LoanRecordCheckerRequest, apiKey, jobId, memberId, companyId string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.ProductCatalogHost + "/product/compliance/loan-record-checker"

	jsonBodyValue, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
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
