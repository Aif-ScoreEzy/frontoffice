package multipleloan

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
	CallMultipleLoan7Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error)
	CallMultipleLoan30Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error)
	CallMultipleLoan90Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error)
}

func (repo *repository) CallMultipleLoan7Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/7-days"

	jsonBodyValue, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-API-Key", apiKey)

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallMultipleLoan30Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/30-days"

	jsonBodyValue, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-API-Key", apiKey)

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallMultipleLoan90Days(request *MultipleLoanRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.ProductCatalogHost + "/product/compliance/multiple-loan/90-days"

	jsonBodyValue, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-API-Key", apiKey)

	client := http.Client{}

	return client.Do(httpRequest)
}
