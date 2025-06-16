package loanrecordchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/internal/httpclient"
	"log"
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

	log.Println("loan record request ====>", httpRequest)

	response, err := repo.Client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}
