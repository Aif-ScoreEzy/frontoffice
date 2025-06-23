package transaction

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repository
}

type Service interface {
	// scoreezy
	GetLogScoreezy() (*AifResponse, int, error)
	GetLogScoreezyByDate(companyId, date string) (*AifResponse, int, error)
	GetLogScoreezyByRangeDate(startDate, endDate, companyId, page string) (*AifResponse, int, error)
	GetLogScoreezyByMonth(companyId, month string) (*AifResponse, int, error)

	// product catalog
	GetLogTransSuccessCount(jobId string) (*getSuccessCountDataResponse, error)
	UpdateLogProCat(transId string, req *UpdateTransRequest) error
}
