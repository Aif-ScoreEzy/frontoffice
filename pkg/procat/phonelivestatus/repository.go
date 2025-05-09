package phonelivestatus

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
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
	q.Add("company_id", filter.CompanyId)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	httpRequest.URL.RawQuery = q.Encode()

	fmt.Println("hittttt", httpRequest)

	client := http.Client{}

	return client.Do(httpRequest)
}
