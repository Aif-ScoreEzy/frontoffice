package livestatus

import (
	"encoding/json"
	"errors"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
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
	CreateJob(data []LiveStatusRequest, userID string, totalData int) (uint, error)
	GetJobs(page, limit, startDate, endDate string) ([]*Job, error)
	GetJobByID(jobID uint) (*Job, error)
	GetJobsTotal(startDate, endDate string) (int64, error)
	GetJobsTotalByRangeDate(startTime, endTime string) (int64, error)
	GetJobDetailsTotalPercentageByRangeDate(startDate, endDate, status string) (int64, error)
	GetJobDetailsPercentageByDataAndRangeDate(startDate, endDate, column, keyword string) (int64, error)
	GetJobDetailsByID(jobID uint) ([]*JobDetail, error)
	GetJobDetailsByRangeDate(startTime, endTime string) ([]*JobDetailQueryResult, error)
	GetJobDetailsWithPagination(page, limit, keyword string, jobID uint) ([]*JobDetailQueryResult, error)
	GetJobDetailsWithPaginationTotal(keyword string, jobID uint) (int64, error)
	GetJobDetailsWithPaginationTotalPercentage(jobID uint, status string) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error)
	GetFailedJobDetails() ([]*JobDetail, error)
	ProcessJobDetails(jobDetail *JobDetail, successRequestTotal int) (int, error)
	CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error)
	UpdateJob(id uint, req *UpdateJobRequest) error
	UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus, status string, data *JSONB) error
	UpdateFailedJobDetail(id uint, sequence int) error
	DeleteJobDetail(id uint) error
	DeleteJob(id uint) error
}

func (svc *service) CreateJob(data []LiveStatusRequest, userID string, totalData int) (uint, error) {
	dataJob := &Job{
		UserID: userID,
		Total:  totalData,
	}

	jobID, err := svc.Repo.CreateJobInTx(dataJob, data)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func (svc *service) GetJobs(page, limit, startDate, endDate string) ([]*Job, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return svc.Repo.GetJobs(intLimit, offset, userID, startTime, endTime)
}

func (svc *service) GetJobsTotalByRangeDate(startDate, endDate string) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}

	return svc.Repo.GetJobsTotalByRangeDate(startTime, endTime)
}

func (svc *service) GetJobByID(jobID uint) (*Job, error) {
	return svc.Repo.GetJobByID(jobID)
}

func (svc *service) GetJobsTotal(startDate, endDate string) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}
	count, err := svc.Repo.GetJobsTotal(startTime, endTime)

	return count, err
}

func (svc *service) GetJobDetailsByID(jobID uint) ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetJobDetailsByJobID(jobID)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsByRangeDate(startDate, endDate string) ([]*JobDetailQueryResult, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return nil, err
	}

	jobDetails, err := svc.Repo.GetJobDetailsByRangeDate(startTime, endTime)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsWithPagination(page, limit, keyword string, jobID uint) ([]*JobDetailQueryResult, error) {
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

func (svc *service) GetJobDetailsWithPaginationTotalPercentage(jobID uint, status string) (int64, error) {
	count, err := svc.Repo.GetJobDetailsByJobIDWithPaginationTotaPercentage(jobID, status)
	return count, err
}

func (svc *service) GetJobDetailsTotalPercentageByRangeDate(startDate, endDate, status string) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}

	count, err := svc.Repo.GetJobDetailsTotalPercentageByStatusAndRangeDate(startTime, endTime, status)
	return count, err
}

func (svc *service) GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsPercentage(column, keyword, jobID)
	return count, err
}

func (svc *service) GetJobDetailsPercentageByDataAndRangeDate(startDate, endDate, column, keyword string) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}

	return svc.Repo.GetJobDetailsPercentageByDataAndRangeDate(startTime, endTime, column, keyword)
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
		updateReq := UpdateJobRequest{
			Total: &successRequestTotal,
		}
		err = svc.UpdateJob(jobDetail.JobID, &updateReq)
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

func (svc *service) UpdateJob(id uint, req *UpdateJobRequest) error {
	data := map[string]interface{}{}

	if req.Total != nil {
		data["success"] = *req.Total
	}

	if req.Status != nil {
		data["status"] = *req.Status
	}

	if req.EndAt != nil {
		data["end_at"] = *req.EndAt
	}

	return svc.Repo.UpdateJob(id, data)
}

func (svc *service) UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus, status string, data *JSONB) error {
	updateJobDetail := map[string]interface{}{}

	updateJobDetail["subscriber_status"] = subcriberStatus
	updateJobDetail["device_status"] = deviceStatus
	updateJobDetail["status"] = status
	updateJobDetail["on_process"] = false
	updateJobDetail["data"] = data

	return svc.Repo.UpdateJobDetail(id, updateJobDetail)
}

func (svc *service) UpdateFailedJobDetail(jobID uint, sequence int) error {
	updateJobDetail := map[string]interface{}{}

	updateJobDetail["on_process"] = true
	updateJobDetail["sequence"] = sequence + 1
	updateJobDetail["status"] = "error"

	return svc.Repo.UpdateJobDetail(jobID, updateJobDetail)
}

func (svc *service) DeleteJobDetail(id uint) error {
	return svc.Repo.DeleteJobDetail(id)
}

func (svc *service) DeleteJob(id uint) error {
	return svc.Repo.DeleteJob(id)
}

func formatTime(startDate, endDate string) (string, string, error) {
	var startTime, endTime string
	layoutPostgreDate := "2006-01-02"
	if startDate != "" {
		err := helper.ParseDate(layoutPostgreDate, startDate)
		if err != nil {
			return "", "", errors.New(constant.InvalidDateFormat)
		}

		startTime = helper.FormatStartTimeForSQL(startDate)

		if endDate == "" {
			endTime = helper.FormatEndTimeForSQL(startDate)
		}
	}

	if endDate != "" {
		err := helper.ParseDate(layoutPostgreDate, endDate)
		if err != nil {
			return "", "", errors.New(constant.InvalidDateFormat)
		}

		endTime = helper.FormatEndTimeForSQL(endDate)
	}

	return startTime, endTime, nil
}
