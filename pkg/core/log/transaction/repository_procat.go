package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"net/http"
	"time"
)

func (repo *repository) CallCreateLogTransAPI(payload *LogTransProCatRequest) error {
	url := fmt.Sprintf("%s/api/core/logging/transaction/product-catalog", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallProcessedLogCount(jobId string) (*getProcessedCountResp, error) {
	url := fmt.Sprintf("%s/api/core/logging/transaction/product-catalog/%s/processed_count", repo.cfg.Env.AifcoreHost, jobId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*getProcessedCountResp](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) CallGetLogTransByJobAPI(jobId, companyId string) ([]*LogTransProductCatalog, error) {
	url := fmt.Sprintf("%s/api/core/logging/transaction/product-catalog/%s", repo.cfg.Env.AifcoreHost, jobId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XCompanyId, companyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[[]*LogTransProductCatalog](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) CallUpdateLogTransAPI(transId string, reqBody map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/core/logging/transaction/product-catalog/%s", repo.cfg.Env.AifcoreHost, transId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[*any](resp)
	if err != nil {
		return err
	}

	return nil
}
