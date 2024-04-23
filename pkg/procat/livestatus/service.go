package livestatus

import (
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"io"
	"strconv"
)

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{Cfg: cfg, Repo: repo}
}

type service struct {
	Cfg  *config.Config
	Repo Repository
}

type Service interface {
	CreateJob(data []LiveStatusRequest, totalData int) (uint, error)
	GetJobs(page, limit string) ([]*Job, error)
	GetJobByID(jobID uint) (*Job, error)
	GetJobsTotal() (int64, error)
	GetJobDetails(jobID uint) ([]*JobDetail, error)
	GetJobDetailsWithPagination(page, limit, keyword string, jobID uint) ([]*JobDetail, error)
	GetJobDetailsWithPaginationTotal(keyword string, jobID uint) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error)
	GetUnprocessedJobDetails() ([]*JobDetail, error)
	ProcessJobDetails(jobDetail *JobDetail, successRequestTotal int) (int, error)
	CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error)
	UpdateJob(id uint, total int) error
	UpdateProcessedJobDetail(id uint) error
	UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus string) error
	UpdateFailedJobDetail(id uint, sequence int) error
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

func (svc *service) GetJobs(page, limit string) ([]*Job, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	return svc.Repo.GetJobs(intLimit, offset)
}

func (svc *service) GetJobByID(jobID uint) (*Job, error) {
	return svc.Repo.GetJobByID(jobID)
}

func (svc *service) GetJobsTotal() (int64, error) {
	count, err := svc.Repo.GetJobsTotal()
	return count, err
}

func (svc *service) GetJobDetails(jobID uint) ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetJobDetailsByJobID(jobID)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsWithPagination(page, limit, keyword string, jobID uint) ([]*JobDetail, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	jobDetails, err := svc.Repo.GetJobDetailsByJobIDWithPagination(intLimit, offset, keyword, jobID)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsWithPaginationTotal(keyword string, jobID uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsByJobIDWithPaginationTotal(keyword, jobID)
	return count, err
}

func (svc *service) GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsPercentage(column, keyword, jobID)
	return count, err
}

func (svc *service) GetUnprocessedJobDetails() ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetUnprocessedJobDetails()
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) ProcessJobDetails(jobDetail *JobDetail, successRequestTotal int) (int, error) {
	apiKey := svc.Cfg.Env.ApiKeyLiveStatus
	jobIDStr := strconv.FormatUint(uint64(jobDetail.JobID), 10)
	request := &LiveStatusRequest{
		PhoneNumber: jobDetail.PhoneNumber,
		TrxID:       jobIDStr,
	}

	liveStatusResponse, err := svc.CreateLiveStatus(request, apiKey)
	if err != nil {
		return 0, err
	}

	// todo: jika status code 200 kirim job detail ke aifcore

	dataMap := liveStatusResponse.Data.(map[string]interface{})
	dataLiveMap := dataMap["live"].(map[string]interface{})
	subscriberStatus := fmt.Sprintf("%v", dataLiveMap["subscriber_status"])
	deviceStatus := fmt.Sprintf("%v", dataLiveMap["device_status"])

	// todo: jika status code 200 maka hapus job detail pada temp tabel. Sampai aifcore menyediakan API untuk get job details, untuk sementara jika status code 200 lakukan update subcriber_status dan device_status pada job detail
	if liveStatusResponse.StatusCode == 200 {
		successRequestTotal += 1
		// err = svc.DeleteJobDetail(jobDetail.ID)
		err = svc.UpdateSucceededJobDetail(jobDetail.ID, subscriberStatus, deviceStatus)
		if err != nil {
			return 0, err
		}

		// todo: jika dari aifcore sudah tersedia api untuk get jobs, hapus program update job
		err = svc.UpdateJob(jobDetail.JobID, successRequestTotal)
		if err != nil {
			return 0, err
		}
	} else {
		_ = svc.UpdateFailedJobDetail(jobDetail.JobID, jobDetail.Sequence)
		if err != nil {
			return 0, err
		}
	}

	return successRequestTotal, nil
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

func (svc *service) UpdateJob(id uint, total int) error {
	return svc.Repo.UpdateJob(id, total)
}

func (svc *service) UpdateProcessedJobDetail(id uint) error {
	request := &UpdateJobDetailRequest{
		OnProcess: true,
	}

	return svc.Repo.UpdateJobDetail(id, request)
}

func (svc *service) UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus string) error {
	request := &UpdateJobDetailRequest{
		SubscriberStatus: subcriberStatus,
		DeviceStatus:     deviceStatus,
	}

	return svc.Repo.UpdateJobDetail(id, request)
}

func (svc *service) UpdateFailedJobDetail(jobID uint, sequence int) error {
	request := &UpdateJobDetailRequest{
		OnProcess: false,
		Sequence:  sequence + 1,
	}

	return svc.Repo.UpdateJobDetail(jobID, request)
}

func (svc *service) DeleteJobDetail(id uint) error {
	return svc.Repo.DeleteJobDetail(id)
}

func (svc *service) DeleteJob(id uint) error {
	return svc.Repo.DeleteJob(id)
}
