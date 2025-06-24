package log

import (
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repository
}

type Service interface {
	CreateProCatJob(req *CreateJobRequest) (*createJobDataResponse, error)
	UpdateJobAPI(jobId string, req *UpdateJobRequest) error
	GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
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
	response, err := svc.repo.CallGetProCatJobAPI(filter)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}

func (svc *service) GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	response, err := svc.repo.CallGetProCatJobDetailAPI(filter)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}
