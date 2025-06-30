package phonelivestatus

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
	"front-office/internal/jsonutil"
	"net/http"
	"time"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient, marshalFn jsonutil.Marshaller) Repository {
	if marshalFn == nil {
		marshalFn = json.Marshal
	}

	return &repository{
		cfg:       cfg,
		client:    client,
		marshalFn: marshalFn,
	}
}

type repository struct {
	cfg       *config.Config
	client    httpclient.HTTPClient
	marshalFn jsonutil.Marshaller
}

type Repository interface {
	CallCreateJobAPI(payload *createJobRequest) (*createJobRespData, error)
	CallGetPhoneLiveStatusJobAPI(filter *phoneLiveStatusFilter) (*jobListRespData, error)
	CallGetJobDetailsAPI(filter *phoneLiveStatusFilter) (*jobDetailRespData, error)
	CallGetAllJobDetailsAPI(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	CallGetJobDetailsByRangeDateAPI(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	CallGetJobsSummary(filter *phoneLiveStatusFilter) (*jobsSummaryRespData, error)
	CallGetProcessedCountAPI(jobId string) (*getSuccessCountRespData, error)
	CallUpdateJob(jobId string, req *updateJobRequest) error
	CallUpdateJobDetail(jobId, jobDetailId string, req *updateJobDetailRequest) error
	CallPhoneLiveStatusAPI(apiKey string, payload *phoneLiveStatusRequest) (*model.ProCatAPIResponse[phoneLiveStatusRespData], error)
}

func (repo *repository) CallCreateJobAPI(payload *createJobRequest) (*createJobRespData, error) {
	url := fmt.Sprintf("%s/api/core/phone-live-status/jobs", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := repo.marshalFn(payload)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, payload.MemberId)
	req.Header.Set(constant.XCompanyId, payload.CompanyId)

	resp, err := repo.client.Do(req)
	fmt.Println("pppppp", req, resp)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*createJobRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetPhoneLiveStatusJobAPI(filter *phoneLiveStatusFilter) (*jobListRespData, error) {
	url := fmt.Sprintf("%s/api/core/phone-live-status/jobs", repo.cfg.Env.AifcoreHost)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
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

	apiResp, err := helper.ParseAifcoreAPIResponse[*jobListRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetJobDetailsAPI(filter *phoneLiveStatusFilter) (*jobDetailRespData, error) {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v/details`, repo.cfg.Env.AifcoreHost, filter.JobId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, filter.MemberId)
	req.Header.Set(constant.XCompanyId, filter.CompanyId)
	req.Header.Set(constant.XTierLevel, filter.TierLevel)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("keyword", filter.Keyword)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*jobDetailRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetAllJobDetailsAPI(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error) {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v`, repo.cfg.Env.AifcoreHost, filter.JobId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, filter.MemberId)
	req.Header.Set(constant.XCompanyId, filter.CompanyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[[]*mstPhoneLiveStatusJobDetail](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetJobDetailsByRangeDateAPI(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error) {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/job-details-by-range-date`, repo.cfg.Env.AifcoreHost)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, filter.MemberId)
	req.Header.Set(constant.XCompanyId, filter.CompanyId)
	req.Header.Set(constant.XTierLevel, filter.TierLevel)

	q := req.URL.Query()
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[[]*mstPhoneLiveStatusJobDetail](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetJobsSummary(filter *phoneLiveStatusFilter) (*jobsSummaryRespData, error) {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs-summary`, repo.cfg.Env.AifcoreHost)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XMemberId, filter.MemberId)
	req.Header.Set(constant.XCompanyId, filter.CompanyId)
	req.Header.Set(constant.XTierLevel, filter.TierLevel)

	q := req.URL.Query()
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*jobsSummaryRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetProcessedCountAPI(jobId string) (*getSuccessCountRespData, error) {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v/success_count`, repo.cfg.Env.AifcoreHost, jobId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*getSuccessCountRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallUpdateJob(jobId string, payload *updateJobRequest) error {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v`, repo.cfg.Env.AifcoreHost, jobId)

	bodyBytes, err := repo.marshalFn(payload)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
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

func (repo *repository) CallUpdateJobDetail(jobId, jobDetailId string, payload *updateJobDetailRequest) error {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v/details/%v`, repo.cfg.Env.AifcoreHost, jobId, jobDetailId)

	bodyBytes, err := repo.marshalFn(payload)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
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

func (repo *repository) CallPhoneLiveStatusAPI(apiKey string, payload *phoneLiveStatusRequest) (*model.ProCatAPIResponse[phoneLiveStatusRespData], error) {
	url := repo.cfg.Env.ProductCatalogHost + "/product/identity/phone-live-status"

	bodyBytes, err := repo.marshalFn(payload)
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

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseProCatAPIResponse[phoneLiveStatusRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}
