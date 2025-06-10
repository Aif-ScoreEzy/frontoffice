package phonelivestatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"io"
	"mime/multipart"
	"net/http"
)

func NewRepository(cfg *config.Config) Repository {
	return &repository{Cfg: cfg}
}

type repository struct {
	Cfg *config.Config
}

type Repository interface {
	CallGetPhoneLiveStatusJobAPI(filter *PhoneLiveStatusFilter) (*http.Response, error)
	CallGetJobDetailsAPI(filter *PhoneLiveStatusFilter) (*http.Response, error)
	CallGetAllJobDetailsAPI(filter *PhoneLiveStatusFilter) (*http.Response, error)
	CallGetJobDetailsByRangeDateAPI(filter *PhoneLiveStatusFilter) (*http.Response, error)
	CallGetJobsSummary(filter *PhoneLiveStatusFilter) (*http.Response, error)
	CallPhoneLiveStatusAPI(memberId, companyId string, request *PhoneLiveStatusRequest) (*http.Response, error)
	CallBulkPhoneLiveStatusAPI(memberId, companyId string, fileHeader *multipart.FileHeader) (*http.Response, error)
}

func (repo *repository) CallGetPhoneLiveStatusJobAPI(filter *PhoneLiveStatusFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/phone-live-status/jobs"

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	q := httpRequest.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	httpRequest.URL.RawQuery = q.Encode()

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallGetJobDetailsAPI(filter *PhoneLiveStatusFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/phone-live-status/job/%v/details`, repo.Cfg.Env.AifcoreHost, filter.JobId)

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	q := httpRequest.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Size)
	q.Add("keyword", filter.Keyword)
	httpRequest.URL.RawQuery = q.Encode()

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallGetAllJobDetailsAPI(filter *PhoneLiveStatusFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/phone-live-status/job/%v`, repo.Cfg.Env.AifcoreHost, filter.JobId)

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallGetJobDetailsByRangeDateAPI(filter *PhoneLiveStatusFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/phone-live-status/job-details-by-range-date`, repo.Cfg.Env.AifcoreHost)

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	q := httpRequest.URL.Query()
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	httpRequest.URL.RawQuery = q.Encode()

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallGetJobsSummary(filter *PhoneLiveStatusFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/phone-live-status/jobs-summary`, repo.Cfg.Env.AifcoreHost)

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", filter.MemberId)
	httpRequest.Header.Set("X-Company-ID", filter.CompanyId)
	httpRequest.Header.Set("X-Tier-Level", filter.TierLevel)

	q := httpRequest.URL.Query()
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	httpRequest.URL.RawQuery = q.Encode()

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallPhoneLiveStatusAPI(memberId, companyId string, request *PhoneLiveStatusRequest) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/phone-live-status/single-search"

	jsonBodyValue, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpRequest.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	httpRequest.Header.Set("X-Member-ID", memberId)
	httpRequest.Header.Set("X-Company-ID", companyId)

	fmt.Println("phone live status reqqq==> ", httpRequest)

	client := http.Client{}

	return client.Do(httpRequest)
}

func (repo *repository) CallBulkPhoneLiveStatusAPI(memberId, companyId string, fileHeader *multipart.FileHeader) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/phone-live-status/bulk-search"

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
