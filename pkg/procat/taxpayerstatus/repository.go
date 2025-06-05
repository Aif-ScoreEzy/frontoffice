package taxpayerstatus

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
	CallTaxPayerStatusAPI(apiKey string, request *taxPayerStatusRequest) (*http.Response, error)
}

func (repo *repository) CallTaxPayerStatusAPI(apiKey string, request *taxPayerStatusRequest) (*http.Response, error) {
	apiUrl := repo.cfg.Env.ProductCatalogHost + "/product/incometax/tax-payer-status"

	jsonBody, err := json.Marshal(request)
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
	httpRequest.Header.Set(constant.XAPIKey, apiKey)

	resp, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}
