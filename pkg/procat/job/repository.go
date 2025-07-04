package job

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
	CallCreateProCatJobAPI(reqBody *CreateJobRequest) (*createJobDataResponse, error)
	CallUpdateJobAPI(jobId string, req map[string]interface{}) error
	CallGetProCatJobAPI(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	CallGetProCatJobDetailAPI(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
}

func (repo *repository) CallCreateProCatJobAPI(reqBody *CreateJobRequest) (*createJobDataResponse, error) {
	url := fmt.Sprintf("%s/api/core/product/jobs", repo.cfg.Env.AifcoreHost)

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
	req.Header.Set(constant.XMemberId, reqBody.MemberId)
	req.Header.Set(constant.XCompanyId, reqBody.CompanyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*createJobDataResponse](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) CallUpdateJobAPI(jobId string, reqBody map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/core/product/jobs/%s", repo.cfg.Env.AifcoreHost, jobId)

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

	_, err = helper.ParseAifcoreAPIResponse[*createJobDataResponse](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallGetProCatJobAPI(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	url := fmt.Sprintf("%s/api/core/product/%s/jobs", repo.cfg.Env.AifcoreHost, filter.ProductSlug)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, filter.MemberId)
	req.Header.Set(constant.XCompanyId, filter.CompanyId)
	req.Header.Set(constant.XTierLevel, filter.TierLevel)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

func (repo *repository) CallGetProCatJobDetailAPI(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	url := fmt.Sprintf("%s/api/core/product/%s/jobs/%s", repo.cfg.Env.AifcoreHost, filter.ProductSlug, filter.JobId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, filter.MemberId)
	req.Header.Set(constant.XCompanyId, filter.CompanyId)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}
