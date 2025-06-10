package transaction

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
	FetchLogTransactions() (*http.Response, error)
	FetchLogTransactionsByDate(companyId, date string) (*http.Response, error)
	FetchLogTransactionsByRangeDate(companyId, startDate, endDate string) (*http.Response, error)
	FetchLogTransactionsByMonth(companyId, month string) (*http.Response, error)
}

func (repo *repository) FetchLogTransactions() (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/list"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) FetchLogTransactionsByDate(companyId, date string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/by"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("date", date)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) FetchLogTransactionsByRangeDate(companyId, startDate, endDate string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/range"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("date_start", startDate)
	q.Add("date_end", endDate)
	q.Add("company_id", companyId)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) FetchLogTransactionsByMonth(companyId, month string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/month"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("month", month)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}
