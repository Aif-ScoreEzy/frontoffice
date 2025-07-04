package transaction

import (
	"front-office/internal/apperror"
)

func (svc *service) GetProcessedLogCount(jobId string) (*getProcessedCountResp, error) {
	result, err := svc.repo.CallProcessedLogCount(jobId)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to get success count")
	}

	return result, nil
}

func (svc *service) UpdateLogProCat(transId string, req *UpdateTransRequest) error {
	data := map[string]interface{}{}

	if req.Success != nil {
		data["success"] = *req.Success
	}

	err := svc.repo.CallUpdateLogTransAPI(transId, data)
	if err != nil {
		return apperror.MapRepoError(err, "failed to update log")
	}

	return nil
}
