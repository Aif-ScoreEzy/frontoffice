package livestatus

import (
	"encoding/json"
	"front-office/app/config"
	"io"
	"log"
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
	GetJobDetailsWithPaginationTotalPercentage(jobID uint) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error)
	GetFailedJobDetails() ([]*JobDetail, error)
	ProcessJobDetails(jobDetail *JobDetail, successRequestTotal int) (int, error)
	CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error)
	UpdateJob(id uint, total int) error
	UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus, status string, data *JSONB) error
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

func (svc *service) GetJobDetailsWithPaginationTotalPercentage(jobID uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsByJobIDWithPaginationTotaPercentage(jobID)
	return count, err
}

func (svc *service) GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsPercentage(column, keyword, jobID)
	return count, err
}

func (svc *service) GetFailedJobDetails() ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetFailedJobDetails()
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

	dataMap, ok := liveStatusResponse.Data.(map[string]interface{})
	if !ok {
		log.Println("Failed to assert Data field as map[string]interface{}")
	}

	dataLiveMap, ok := dataMap["live"].(map[string]interface{})
	if !ok {
		log.Println("Failed to assert live field within Data as map[string]interface{}")
	}

	subscriberStatus, ok := dataLiveMap["subscriber_status"].(string)
	if !ok {
		log.Println("Failed to assert subscriber_status field as string")
	}

	deviceStatus, ok := dataLiveMap["device_status"].(string)
	if !ok {
		log.Println("Failed to assert device_status field as string")
	}

	var errorCode int
	if errors, ok := dataMap["errors"].([]interface{}); ok {
		for _, err := range errors {
			if errMap, ok := err.(map[string]interface{}); ok {
				if code, ok := errMap["code"].(float64); ok {
					errorCode = int(code)
				} else {
					log.Println("Error: 'code' field is not a number")
				}
			}
		}
	}

	data := &JSONB{}
	responseBodyByte, err := json.Marshal(liveStatusResponse.Data)
	if err == nil {
		(*data).Scan(responseBodyByte)
	}

	// todo: jika status code 200 maka hapus job detail pada temp tabel. Sampai aifcore menyediakan API untuk get job details, untuk sementara jika status code 200 lakukan update subcriber_status dan device_status pada job detail
	if liveStatusResponse.StatusCode == 200 {
		// todo: pastikan errors bukan kode 6001, update kolom status "success", jika errors code 6001 update status "fail", hanya status "error" yg diulang
		successRequestTotal += 1
		// err = svc.DeleteJobDetail(jobDetail.ID)
		if errorCode == -60001 {
			err = svc.UpdateSucceededJobDetail(jobDetail.ID, subscriberStatus, deviceStatus, "fail", data)
			if err != nil {
				return 0, err
			}
		} else {
			err = svc.UpdateSucceededJobDetail(jobDetail.ID, subscriberStatus, deviceStatus, "success", data)
			if err != nil {
				return 0, err
			}
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

func (svc *service) UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus, status string, data *JSONB) error {
	request := &UpdateJobDetailRequest{
		OnProcess:        false,
		SubscriberStatus: subcriberStatus,
		DeviceStatus:     deviceStatus,
		Status:           status,
		Data:             data,
	}

	return svc.Repo.UpdateJobDetail(id, request)
}

func (svc *service) UpdateFailedJobDetail(jobID uint, sequence int) error {
	request := &UpdateJobDetailRequest{
		OnProcess: true,
		Sequence:  sequence + 1,
		Status:    "error",
	}

	return svc.Repo.UpdateJobDetail(jobID, request)
}

func (svc *service) DeleteJobDetail(id uint) error {
	return svc.Repo.DeleteJobDetail(id)
}

func (svc *service) DeleteJob(id uint) error {
	return svc.Repo.DeleteJob(id)
}
