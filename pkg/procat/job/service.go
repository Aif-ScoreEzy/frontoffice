package job

import (
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/transaction"
	"time"
)

func NewService(repo Repository, transactionRepo transaction.Repository) Service {
	return &service{
		repo,
		transactionRepo,
	}
}

type service struct {
	repo            Repository
	transactionRepo transaction.Repository
}

type Service interface {
	CreateProCatJob(req *CreateJobRequest) (*createJobDataResponse, error)
	UpdateJobAPI(jobId string, req *UpdateJobRequest) error
	GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	FinalizeJob(jobIdStr string, transactionId string) error
	FinalizeFailedJob(jobIdStr string) error
}

func (svc *service) CreateProCatJob(req *CreateJobRequest) (*createJobDataResponse, error) {
	result, err := svc.repo.CallCreateProCatJobAPI(req)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to create job")
	}

	return result, nil
}

func (svc *service) UpdateJobAPI(jobId string, req *UpdateJobRequest) error {
	data := map[string]interface{}{}

	if req.SuccessCount != nil {
		data["success_count"] = *req.SuccessCount
	}

	if req.Status != nil {
		data["status"] = *req.Status
	}

	if req.EndAt != nil {
		data["end_at"] = *req.EndAt
	}

	err := svc.repo.CallUpdateJobAPI(jobId, data)
	if err != nil {
		return apperror.MapRepoError(err, "failed to update job")
	}

	return nil
}

func (svc *service) GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	result, err := svc.repo.CallGetProCatJobAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch jobs")
	}

	return result, nil
}

func (svc *service) GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	result, err := svc.repo.CallGetProCatJobDetailAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch job detail")
	}

	return result, nil
}

func (svc *service) FinalizeJob(jobIdStr string, transactionId string) error {
	if err := svc.transactionRepo.CallUpdateLogTransAPI(transactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update transaction log")
	}

	count, err := svc.transactionRepo.CallProcessedLogCount(jobIdStr)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get success count")
	}

	if err := svc.repo.CallUpdateJobAPI(jobIdStr, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.ProcessedCount)),
		"status":        helper.StringPtr(constant.JobStatusDone),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update job status")
	}

	return nil
}

func (svc *service) FinalizeFailedJob(jobIdStr string) error {
	count, err := svc.transactionRepo.CallProcessedLogCount(jobIdStr)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get processed count request")
	}

	if err := svc.repo.CallUpdateJobAPI(jobIdStr, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.ProcessedCount)),
		"status":        helper.StringPtr(constant.JobStatusFailed),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update job status")
	}

	return nil
}
