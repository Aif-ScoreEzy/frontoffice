package loanrecordchecker

import (
	"bytes"
	"encoding/json"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
)

func NewRepository(cfg *config.Config) Repository {
	return &repository{Cfg: cfg}
}

type repository struct {
	Cfg *config.Config
}

type Repository interface {
	CallLoanRecordChecker(request *LoanRecordCheckerRequest, apiKey string) (*http.Response, error)
}

func (repo *repository) CallLoanRecordChecker(request *LoanRecordCheckerRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/loan-record-checker"

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

	client := http.Client{}

	return client.Do(httpRequest)
}
