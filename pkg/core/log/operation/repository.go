package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/httpclient"
	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{cfg, client}
}

type repository struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type Repository interface {
	CallGetLogsOperationAPI(filter *LogOperationFilter) (*model.AifcoreAPIResponse[any], error)
	CallGetLogsByRangeAPI(filter *LogRangeFilter) (*model.AifcoreAPIResponse[any], error)
	AddLogOperation(req *AddLogRequest) error
}

func (repo *repository) CallGetLogsOperationAPI(filter *LogOperationFilter) (*model.AifcoreAPIResponse[any], error) {
	url := fmt.Sprintf("%s/api/core/logging/operation/list/%s", repo.cfg.Env.AifcoreHost, filter.CompanyId)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("name", filter.Name)
	q.Add("role", filter.Role)
	q.Add("event", filter.Event)
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

func (repo *repository) CallGetLogsByRangeAPI(filter *LogRangeFilter) (*model.AifcoreAPIResponse[any], error) {
	url := fmt.Sprintf("%s/api/core/logging/operation/range", repo.cfg.Env.AifcoreHost)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("company_id", filter.CompanyId)
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

func (repo *repository) AddLogOperation(reqBody *AddLogRequest) error {
	url := fmt.Sprintf("%s/api/core/logging/operation", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}
