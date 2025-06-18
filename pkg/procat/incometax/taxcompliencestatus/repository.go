package taxcompliancestatus

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
	CallTaxComplianceStatusAPI(apiKey, jobId string, request *taxComplianceStatusRequest) (*http.Response, error)
}

func (repo *repository) CallTaxComplianceStatusAPI(apiKey, jobId string, request *taxComplianceStatusRequest) (*http.Response, error) {
	apiUrl := repo.cfg.Env.ProductCatalogHost + "/product/incometax/tax-compliance-status"

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

	q := httpRequest.URL.Query()
	q.Add("job_id", jobId)
	httpRequest.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}
