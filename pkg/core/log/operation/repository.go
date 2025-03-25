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
	FetchLogOperations(companyId string, filter *GetLogOperationFilter) (*http.Response, error)
}

func (repo *repository) FetchLogOperations(companyId string, filter *GetLogOperationFilter) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/operation/list/" + companyId

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("role", filter.Role)
	q.Add("event", filter.Event)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}
