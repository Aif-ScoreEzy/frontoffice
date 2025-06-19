package transaction

import (
	"front-office/common/constant"
	"net/http"
)

func (repo *repository) CallLogScoreezyAPI() (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/list"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) CallLogScoreezyByDateAPI(companyId, date string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/by"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("date", date)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) CallLogScoreezyByRangeDateAPI(companyId, startDate, endDate string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/range"

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

func (repo *repository) CallLogScoreezyByMonthAPI(companyId, month string) (*http.Response, error) {
	apiUrl := repo.cfg.Env.AifcoreHost + "/api/core/logging/transaction/scoreezy/month"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("month", month)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}
