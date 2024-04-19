package livestatus

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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
	ProcessBatchJobDetails(apiKey string, jobID uint, batch []*JobDetail) ([]*LiveStatusResponse, error)
	ProcessJobDetails(apiKey string, jobID uint, jobDetails []*JobDetail, batchSize int) ([]*LiveStatusResponse, error)
	CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error)
	DeleteJobDetail(id uint) error
	DeleteJob(id uint) error
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

func (svc *service) ProcessBatchJobDetails(apiKey string, jobID uint, batch []*JobDetail) ([]*LiveStatusResponse, error) {
	var liveStatusResponse *LiveStatusResponse
	var liveStatusResponses []*LiveStatusResponse
	// var err error
	jobIDStr := strconv.FormatUint(uint64(jobID), 10)

	for _, jobDetail := range batch {
		fmt.Printf("Processing Job Detail ID: %d\n", jobDetail.ID)
		request := &LiveStatusRequest{
			PhoneNumber: jobDetail.PhoneNumber,
			TrxID:       jobIDStr,
		}

		liveStatusResponse, _ = svc.CreateLiveStatus(request, apiKey)
		// if err != nil {
		// 	return nil, err
		// }

		liveStatusResponses = append(liveStatusResponses, liveStatusResponse)

		_ = svc.DeleteJobDetail(jobDetail.ID)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return liveStatusResponses, nil
}

func (svc *service) ProcessJobDetails(apiKey string, jobID uint, jobDetails []*JobDetail, batchSize int) ([]*LiveStatusResponse, error) {
	var liveStatusResponses []*LiveStatusResponse
	var err error
	numJobDetails := len(jobDetails)
	for i := 0; i < numJobDetails; i += batchSize {
		end := i + batchSize
		if end > numJobDetails {
			end = numJobDetails
		}

		batch := jobDetails[i:end]

		fmt.Printf("Processing batch %d to %d\n", i, end)
		liveStatusResponses, err = svc.ProcessBatchJobDetails(apiKey, jobID, batch)
		if err != nil {
			return nil, err
		}
	}

	return liveStatusResponses, nil
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

func (svc *service) DeleteJob(id uint) error {
	return svc.Repo.DeleteJob(id)
}