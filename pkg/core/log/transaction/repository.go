package transaction

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{
		cfg:    cfg,
		client: client,
	}
}

type repository struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type Repository interface {
	// scoreezy
	CallLogScoreezyAPI() (*http.Response, error)
	CallLogScoreezyByDateAPI(companyId, date string) (*http.Response, error)
	CallLogScoreezyByRangeDateAPI(companyId, startDate, endDate string) (*http.Response, error)
	CallLogScoreezyByMonthAPI(companyId, month string) (*http.Response, error)

	// product catalog
	CallLogTransSuccessCountAPI(jobId string) (*getSuccessCountDataResponse, error)
	CallUpdateLogTransAPI(transId string, req map[string]interface{}) error
}
