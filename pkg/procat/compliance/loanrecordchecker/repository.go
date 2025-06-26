package loanrecordchecker

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
	CallLoanRecordCheckerAPI(apiKey, jobId, memberId, companyId string, request *loanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error)
}

func (repo *repository) CallLoanRecordCheckerAPI(apiKey, jobId, memberId, companyId string, reqBody *loanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error) {
	url := fmt.Sprintf("%s/product/compliance/loan-record-checker", repo.cfg.Env.ProductCatalogHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XAPIKey, apiKey)
	req.Header.Set(constant.XMemberId, memberId)
	req.Header.Set(constant.XCompanyId, companyId)

	q := req.URL.Query()
	q.Add("job_id", jobId)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	return helper.ParseProCatAPIResponse[dataLoanRecord](resp)
}
