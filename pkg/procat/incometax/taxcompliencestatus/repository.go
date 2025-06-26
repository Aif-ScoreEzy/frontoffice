package taxcompliancestatus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
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
	CallTaxComplianceStatusAPI(apiKey, jobId string, request *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceDataResponse], error)
}

func (repo *repository) CallTaxComplianceStatusAPI(apiKey, jobId string, reqBody *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceDataResponse], error) {
	url := fmt.Sprintf("%s/product/incometax/tax-compliance-status", repo.cfg.Env.ProductCatalogHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XAPIKey, apiKey)

	q := req.URL.Query()
	q.Add("job_id", jobId)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseProCatAPIResponse[taxComplianceDataResponse](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, err
}
