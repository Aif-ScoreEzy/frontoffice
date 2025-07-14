package transaction

import (
	"encoding/json"
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/internal/jsonutil"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient, marshalFn jsonutil.Marshaller) Repository {
	if marshalFn == nil {
		marshalFn = json.Marshal
	}

	return &repository{
		cfg:       cfg,
		client:    client,
		marshalFn: marshalFn,
	}
}

type repository struct {
	cfg       *config.Config
	client    httpclient.HTTPClient
	marshalFn jsonutil.Marshaller
}

type Repository interface {
	// scoreezy
	GetLogsScoreezyAPI() ([]*LogTransScoreezy, error)
	GetLogsScoreezyByDateAPI(companyId, date string) ([]*LogTransScoreezy, error)
	GetLogsScoreezyByDateRangeAPI(companyId, startDate, endDate string) ([]*LogTransScoreezy, error)
	CallScoreezyLogsByMonthAPI(companyId, month string) ([]*LogTransScoreezy, error)

	// product catalog
	CreateLogTransAPI(req *LogTransProCatRequest) error
	GetLogTransByJobIdAPI(jobId, companyId string) ([]*LogTransProductCatalog, error)
	ProcessedLogCountAPI(jobId string) (*getProcessedCountResp, error)
	UpdateLogTransAPI(transId string, req map[string]interface{}) error
}
