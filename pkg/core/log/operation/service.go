package operation

import (
	"front-office/internal/apperror"
)

func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

type Service interface {
	GetLogsOperation(filter *LogOperationFilter) ([]*LogOperation, error)
	GetLogsByRange(filter *LogRangeFilter) ([]*LogOperation, error)
	AddLogOperation(req *AddLogRequest) error
}

func (svc *service) GetLogsOperation(filter *LogOperationFilter) ([]*LogOperation, error) {
	logs, err := svc.repo.CallGetLogsOperationAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch log operations")
	}

	return logs, nil
}

func (svc *service) GetLogsByRange(filter *LogRangeFilter) ([]*LogOperation, error) {
	logs, err := svc.repo.CallGetLogsByRangeAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch log operations")
	}

	return logs, nil
}

func (svc *service) AddLogOperation(req *AddLogRequest) error {
	if err := svc.repo.AddLogOperation(req); err != nil {
		return apperror.MapRepoError(err, "failed to create log")
	}

	return nil
}
