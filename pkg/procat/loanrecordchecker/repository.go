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
		Cfg:    cfg,
		Client: client,
	}
}

type repository struct {
	Cfg    *config.Config
	Client httpclient.HTTPClient
}

type Repository interface {
	CallLoanRecordCheckerAPI(request *LoanRecordCheckerRequest, apiKey, memberId, companyId string) (*http.Response, error)
	CallGetLoanRecordCheckerJobAPI(filter *loanRecordCheckerFilter) (*http.Response, error)
}

func (repo *repository) CallLoanRecordCheckerAPI(request *LoanRecordCheckerRequest, apiKey, memberId, companyId string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/product/compliance/loan-record-checker"

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

	response, err := repo.Client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}

func (repo *repository) CallGetLoanRecordCheckerJobAPI(filter *loanRecordCheckerFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/job/by-product/" + filter.ProductSlug

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
