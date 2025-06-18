package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/common/constant"
	"net/http"
	"time"
)

func (repo *repository) CallLogTransSuccessCountAPI(jobId string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/transaction/product-catalog/" + jobId

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

func (repo *repository) CallUpdateLogTransAPI(transId string, req map[string]interface{}) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/transaction/product-catalog/" + transId

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPut, apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	fmt.Println("update log", httpRequest, resp)

	return resp, nil
}
