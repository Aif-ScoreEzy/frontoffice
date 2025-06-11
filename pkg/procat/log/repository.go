package log

import (
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/internal/httpclient"
	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{
		Cfg:    cfg,
		Client: client,
	}
}

type repository struct {
	Cfg    *config.Config
	Client httpclient.HTTPClient
}

type Repository interface {
	CallProCatJobAPI(filter *logFilter) (*http.Response, error)
	CallGetProCatJobDetailAPI(filter *logFilter) (*http.Response, error)
}

func (repo *repository) CallProCatJobAPI(filter *logFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/product/" + filter.ProductSlug + "/jobs"

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

func (repo *repository) CallGetProCatJobDetailAPI(filter *logFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/product/" + filter.ProductSlug + "/jobs/" + filter.JobId

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
