package livestatus

import (
	"encoding/json"
	"io"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateJob(data []LiveStatusRequest, totalData int) (uint, error)
	GetJobDetails(jobID uint) ([]*JobDetail, error)
	CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error)
	DeleteJobDetail(id uint) error
}

func (svc *service) CreateJob(data []LiveStatusRequest, totalData int) (uint, error) {
	dataJob := &Job{
		Total: totalData,
	}

	jobID, err := svc.Repo.CreateJobInTx(dataJob, data)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func (svc *service) GetJobDetails(jobID uint) ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetJobDetailsByJobID(jobID)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error) {
	response, err := svc.Repo.CallLiveStatus(liveStatusRequest, apiKey)
	if err != nil {
		return nil, err
	}

	dataBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var liveStatusResponse *LiveStatusResponse
	json.Unmarshal(dataBytes, &liveStatusResponse)
	liveStatusResponse.StatusCode = response.StatusCode

	return liveStatusResponse, nil
}

func (svc *service) DeleteJobDetail(id uint) error {
	return svc.Repo.DeleteJobDetail(id)
}
