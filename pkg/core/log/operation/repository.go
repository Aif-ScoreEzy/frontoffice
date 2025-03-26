package operation

import (
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
)

func NewRepository(cfg *config.Config) Repository {
	return &repository{
		Cfg: cfg,
	}
}

type repository struct {
	Cfg *config.Config
}

type Repository interface {
	FetchLogOperations(filter *LogOperationFilter) (*http.Response, error)
	FetchByRange(filter *LogRangeFilter) (*http.Response, error)
}

func (repo *repository) FetchLogOperations(filter *LogOperationFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/operation/list/" + filter.CompanyId

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("role", filter.Role)
	q.Add("event", filter.Event)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) FetchByRange(filter *LogRangeFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/operation/range"

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", filter.CompanyId)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}
