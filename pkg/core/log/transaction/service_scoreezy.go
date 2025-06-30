package transaction

import (
	"front-office/common/constant"
	"front-office/internal/apperror"
)

func (svc *service) GetScoreezyLogs() ([]*scoreezyLogResponse, error) {
	logs, err := svc.repo.CallScoreezyLogsAPI()
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedFetchLogs)
	}

	result := make([]*scoreezyLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, &scoreezyLogResponse{
			Name:      log.Member.Name,
			Grade:     log.Grade,
			CreatedAt: log.CreatedAt,
		})
	}

	return result, nil
}

func (svc *service) GetScoreezyLogsByDate(companyId, date string) ([]*scoreezyLogResponse, error) {
	logs, err := svc.repo.CallScoreezyLogsByDateAPI(companyId, date)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedFetchLogs)
	}

	result := make([]*scoreezyLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, &scoreezyLogResponse{
			Name:      log.Member.Name,
			Grade:     log.Grade,
			CreatedAt: log.CreatedAt,
		})
	}

	return result, nil
}

func (svc *service) GetScoreezyLogsByRangeDate(startDate, endDate, companyId, page string) ([]*scoreezyLogResponse, error) {
	logs, err := svc.repo.CallScoreezyLogsByRangeDateAPI(companyId, startDate, endDate)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedFetchLogs)
	}

	result := make([]*scoreezyLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, &scoreezyLogResponse{
			Name:      log.Member.Name,
			Grade:     log.Grade,
			CreatedAt: log.CreatedAt,
		})
	}

	return result, nil
}

func (svc *service) GetScoreezyLogsByMonth(companyId, month string) ([]*scoreezyLogResponse, error) {
	logs, err := svc.repo.CallScoreezyLogsByMonthAPI(companyId, month)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedFetchLogs)
	}

	result := make([]*scoreezyLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, &scoreezyLogResponse{
			Name:      log.Member.Name,
			Grade:     log.Grade,
			CreatedAt: log.CreatedAt,
		})
	}

	return result, nil
}
