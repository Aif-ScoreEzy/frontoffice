package loanrecordchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
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
	CallLoanRecordCheckerAPI(apiKey, jobId, memberId, companyId string, request *LoanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error)
}

func (repo *repository) CallLoanRecordCheckerAPI(apiKey, jobId, memberId, companyId string, reqBody *LoanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error) {
	url := fmt.Sprintf("%s/product/compliance/loan-record-checker", repo.cfg.Env.ProductCatalogHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("X-Member-ID", memberId)
	req.Header.Set("X-Company-ID", companyId)

	q := req.URL.Query()
	q.Add("job_id", jobId)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	return helper.ParseProCatAPIResponse[dataLoanRecord](resp)
}
