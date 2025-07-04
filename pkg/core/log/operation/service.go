package operation

import (
	"front-office/common/constant"
	"front-office/internal/apperror"
)

func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

type Service interface {
	GetLogsOperation(filter *LogOperationFilter) (*logOperationAPIResponse, error)
	GetLogsByRange(filter *LogRangeFilter) (*logOperationAPIResponse, error)
	AddLogOperation(req *AddLogRequest) error
}

func (svc *service) GetLogsOperation(filter *LogOperationFilter) (*logOperationAPIResponse, error) {
	result, err := svc.repo.CallGetLogsOperationAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch log operations")
	}

	response := &logOperationAPIResponse{
		Message: constant.SucceedGetLogTrans,
		Success: result.Success,
		Data: &logOperationRespData{
			Logs: result.Data,
		},
		Meta: *result.Meta,
	}

	return response, nil
}

func (svc *service) GetLogsByRange(filter *LogRangeFilter) (*logOperationAPIResponse, error) {
	result, err := svc.repo.CallGetLogsByRangeAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch log operations")
	}

	response := &logOperationAPIResponse{
		Message: constant.SucceedGetLogTrans,
		Success: result.Success,
		Data: &logOperationRespData{
			Logs: result.Data,
		},
		Meta: *result.Meta,
	}

	return response, nil
}

func (svc *service) AddLogOperation(req *AddLogRequest) error {
	if err := svc.repo.AddLogOperation(req); err != nil {
		return apperror.MapRepoError(err, "failed to create log")
	}

	return nil
}
