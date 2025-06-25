package transaction

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
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
	CallScoreezyLogsAPI() ([]*LogTransScoreezy, error)
	CallScoreezyLogsByDateAPI(companyId, date string) ([]*LogTransScoreezy, error)
	CallScoreezyLogsByRangeDateAPI(companyId, startDate, endDate string) ([]*LogTransScoreezy, error)
	CallScoreezyLogsByMonthAPI(companyId, month string) ([]*LogTransScoreezy, error)

	// product catalog
	CallLogTransSuccessCountAPI(jobId string) (*getSuccessCountDataResponse, error)
	CallUpdateLogTransAPI(transId string, req map[string]interface{}) error
}
