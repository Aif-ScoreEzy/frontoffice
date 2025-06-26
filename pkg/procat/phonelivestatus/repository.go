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
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{cfg, client}
}

type repository struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type Repository interface {
	CallCreateJobAPI(memberId, companyId string, req *createJobRequest) (*createJobResponseData, error)
	CallGetPhoneLiveStatusJobAPI(filter *phoneLiveStatusFilter) (*jobListRespData, error)
	CallGetJobDetailsAPI(filter *phoneLiveStatusFilter) (*jobDetailRespData, error)
	CallGetAllJobDetailsAPI(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	CallGetJobDetailsByRangeDateAPI(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	CallGetJobsSummary(filter *phoneLiveStatusFilter) (*jobsSummaryRespData, error)
	CallGetProcessedCount(jobId string) (*getSuccessCountRespData, error)
	CallUpdateJob(jobId string, req *updateJobRequest) error
	CallUpdateJobDetail(jobId, jobDetailId string, req *updateJobDetailRequest) error
	CallPhoneLiveStatusAPI(apiKey string, reqBody *phoneLiveStatusRequest) (*model.ProCatAPIResponse[phoneLiveStatusRespData], error)
	CallBulkPhoneLiveStatusAPI(memberId, companyId string, fileHeader *multipart.FileHeader) (*http.Response, error)
}

func (repo *repository) CallCreateJobAPI(memberId, companyId string, reqBody *createJobRequest) (*createJobResponseData, error) {
	url := fmt.Sprintf("%s/api/core/phone-live-status/jobs", repo.cfg.Env.AifcoreHost)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set("X-Member-ID", memberId)
	req.Header.Set("X-Company-ID", companyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*createJobResponseData](resp)
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
	req.Header.Set("X-Member-ID", filter.MemberId)
	req.Header.Set("X-Company-ID", filter.CompanyId)
	req.Header.Set("X-Tier-Level", filter.TierLevel)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
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
	req.Header.Set("X-Member-ID", filter.MemberId)
	req.Header.Set("X-Company-ID", filter.CompanyId)
	req.Header.Set("X-Tier-Level", filter.TierLevel)

	q := req.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("keyword", filter.Keyword)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
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
	req.Header.Set("X-Member-ID", filter.MemberId)
	req.Header.Set("X-Company-ID", filter.CompanyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
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
	req.Header.Set("X-Member-ID", filter.MemberId)
	req.Header.Set("X-Company-ID", filter.CompanyId)
	req.Header.Set("X-Tier-Level", filter.TierLevel)

	q := req.URL.Query()
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
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
	req.Header.Set("X-Member-ID", filter.MemberId)
	req.Header.Set("X-Company-ID", filter.CompanyId)
	req.Header.Set("X-Tier-Level", filter.TierLevel)

	q := req.URL.Query()
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*jobsSummaryRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallGetProcessedCount(jobId string) (*getSuccessCountRespData, error) {
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
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*getSuccessCountRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, err
}

func (repo *repository) CallUpdateJob(jobId string, reqBody *updateJobRequest) error {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v`, repo.cfg.Env.AifcoreHost, jobId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
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
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallUpdateJobDetail(jobId, jobDetailId string, reqBody *updateJobDetailRequest) error {
	url := fmt.Sprintf(`%v/api/core/phone-live-status/jobs/%v/details/%v`, repo.cfg.Env.AifcoreHost, jobId, jobDetailId)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
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
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallPhoneLiveStatusAPI(apiKey string, reqBody *phoneLiveStatusRequest) (*model.ProCatAPIResponse[phoneLiveStatusRespData], error) {
	url := repo.cfg.Env.ProductCatalogHost + "/product/identity/phone-live-status"

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set("X-API-KEY", apiKey)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseProCatAPIResponse[phoneLiveStatusRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

func (repo *repository) CallBulkPhoneLiveStatusAPI(memberId, companyId string, fileHeader *multipart.FileHeader) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/phone-live-status/bulk-search"

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpRequest.Header.Set(constant.HeaderContentType, writer.FormDataContentType())
	httpRequest.Header.Set("X-Member-ID", memberId)
	httpRequest.Header.Set("X-Company-ID", companyId)

	client := http.Client{}

	return client.Do(httpRequest)
}
