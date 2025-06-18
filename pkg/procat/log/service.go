package log

import (
	"front-office/common/model"
	"front-office/helper"
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
	CreateProCatJob(req *CreateJobRequest) (*model.AifcoreAPIResponse[createJobDataResponse], error)
	UpdateJobAPI(jobId string, req *UpdateJobRequest) (*model.AifcoreAPIResponse[any], error)
	GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
}

func (svc *service) CreateProCatJob(req *CreateJobRequest) (*model.AifcoreAPIResponse[createJobDataResponse], error) {
	response, err := svc.repo.CallCreateProCatJobAPI(req)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[createJobDataResponse](response)
}

func (svc *service) UpdateJobAPI(jobId string, req *UpdateJobRequest) (*model.AifcoreAPIResponse[any], error) {
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

	response, err := svc.repo.CallUpdateJobAPI(jobId, data)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
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
