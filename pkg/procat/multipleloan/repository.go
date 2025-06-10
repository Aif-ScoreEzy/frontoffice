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
		Cfg:    cfg,
		Client: client,
	}
}

type repository struct {
	Cfg    *config.Config
	Client httpclient.HTTPClient
}

type Repository interface {
	CallMultipleLoan7Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error)
	CallMultipleLoan30Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error)
	CallMultipleLoan90Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error)
}

func (repo *repository) CallMultipleLoan7Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/7-days"

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

	response, err := repo.Client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}

func (repo *repository) CallMultipleLoan30Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/30-days"

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

	response, err := repo.Client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}

func (repo *repository) CallMultipleLoan90Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/90-days"

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

	response, err := repo.Client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return response, nil
}
