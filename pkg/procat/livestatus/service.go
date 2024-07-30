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
	CreateJob(data []LiveStatusRequest, userId, companyId string, totalData int) (uint, error)
	GetJobs(page, limit, userId, companyId, startDate, endDate string, tierLevel uint) ([]*Job, error)
	GetJobById(jobId uint) (*Job, error)
	GetJobByIdAndUserId(jobId, tierLevel uint, userId, companyId string) (*Job, error)
	GetJobsTotal(userId, companyId, startDate, endDate string, tierLevel uint) (int64, error)
	GetJobsTotalByRangeDate(userId, companyId, startDate, endDate string, tierLevel uint) (int64, error)
	GetJobDetailsTotalPercentageByRangeDate(userId, companyId, startDate, endDate, status string, tierLevel uint) (int64, error)
	GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startDate, endDate, column, keyword string, tierLevel uint) (int64, error)
	GetJobDetailsByJobId(jobId uint) ([]*JobDetail, error)
	GetJobDetailsByRangeDate(userId, companyId, startTime, endTime string, tierLevel uint) ([]*JobDetailQueryResult, error)
	GetJobDetailsWithPagination(page, limit, keyword string, jobId uint) ([]*JobDetailQueryResult, error)
	GetJobDetailsWithPaginationTotal(keyword string, jobId uint) (int64, error)
	GetJobDetailsWithPaginationTotalPercentage(jobId uint, status string) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobId uint) (int64, error)
	GetFailedJobDetails(jobId uint) ([]*JobDetail, error)
	ProcessJobDetails(jobDetail *JobDetail) error
	CreateLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*LiveStatusResponse, error)
	UpdateJob(id uint, req *UpdateJobRequest) error
	UpdateSucceededJobDetail(id uint, subcriberStatus, deviceStatus, status string, data *JSONB) error
	UpdateFailedJobDetail(id uint, sequence int) error
	UpdateInvalidJobDetail(id uint, errMessage string) error
	DeleteJobDetail(id uint) error
	DeleteJob(id uint) error
	GetJobDetailsByJobIdExport(jobId uint) ([]*JobDetailQueryResult, error)
	GetJobWithIncompleteStatus() ([]uint, error)
	GetOnProcessJobDetails(jobId uint, onProcess bool) ([]uint, error)
	CountOnProcessJobDetails(jobId uint, onProcess bool) (int, error)
}

func (svc *service) CreateJob(data []LiveStatusRequest, userId, companyId string, totalData int) (uint, error) {
	dataJob := &Job{
		UserId:    userId,
		CompanyId: companyId,
		Total:     totalData,
	}

	jobId, err := svc.Repo.CreateJobInTx(userId, companyId, dataJob, data)
	if err != nil {
		return 0, err
	}

	return jobId, nil
}

func (svc *service) GetJobs(page, limit, userId, companyId, startDate, endDate string, tierLevel uint) ([]*Job, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return svc.Repo.GetJobs(intLimit, offset, tierLevel, userId, companyId, startTime, endTime)
}

func (svc *service) GetJobsTotalByRangeDate(userId, companyId, startDate, endDate string, tierLevel uint) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}

	return svc.Repo.GetJobsTotalByRangeDate(userId, companyId, startTime, endTime, tierLevel)
}

func (svc *service) GetJobById(jobId uint) (*Job, error) {
	return svc.Repo.GetJobById(jobId)
}

func (svc *service) GetJobByIdAndUserId(jobId, tierLevel uint, userId, companyId string) (*Job, error) {
	return svc.Repo.GetJobByIdAndUserId(jobId, tierLevel, userId, companyId)
}

func (svc *service) GetJobsTotal(userId, companyId, startDate, endDate string, tierLevel uint) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}
	count, err := svc.Repo.GetJobsTotal(userId, companyId, startTime, endTime, tierLevel)

	return count, err
}

func (svc *service) GetJobDetailsByJobId(jobId uint) ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetJobDetailsByJobId(jobId)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsByRangeDate(userId, companyId, startDate, endDate string, tierLevel uint) ([]*JobDetailQueryResult, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return nil, err
	}

	jobDetails, err := svc.Repo.GetJobDetailsByRangeDate(userId, companyId, startTime, endTime, tierLevel)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsWithPagination(page, limit, keyword string, jobId uint) ([]*JobDetailQueryResult, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	jobDetails, err := svc.Repo.GetJobDetailsByJobIdWithPagination(intLimit, offset, keyword, jobId)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobDetailsWithPaginationTotal(keyword string, jobId uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsByJobIdWithPaginationTotal(keyword, jobId)
	return count, err
}

func (svc *service) GetJobDetailsWithPaginationTotalPercentage(jobId uint, status string) (int64, error) {
	count, err := svc.Repo.GetJobDetailsByJobIdWithPaginationTotaPercentage(jobId, status)
	return count, err
}

func (svc *service) GetJobDetailsTotalPercentageByRangeDate(userId, companyId, startDate, endDate, status string, tierLevel uint) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}

	count, err := svc.Repo.GetJobDetailsTotalPercentageByStatusAndRangeDate(userId, companyId, startTime, endTime, status, tierLevel)
	return count, err
}

func (svc *service) GetJobDetailsPercentage(column, keyword string, jobId uint) (int64, error) {
	count, err := svc.Repo.GetJobDetailsPercentage(column, keyword, jobId)
	return count, err
}

func (svc *service) GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startDate, endDate, column, keyword string, tierLevel uint) (int64, error) {
	startTime, endTime, err := formatTime(startDate, endDate)
	if err != nil {
		return 0, err
	}

	return svc.Repo.GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startTime, endTime, column, keyword, tierLevel)
}

func (svc *service) GetFailedJobDetails(jobId uint) ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetFailedJobDetails(jobId)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) ProcessJobDetails(jobDetail *JobDetail) error {
	apiKey := svc.Cfg.Env.ApiKeyLiveStatus
	jobIdStr := strconv.FormatUint(uint64(jobDetail.JobId), 10)
	request := &LiveStatusRequest{
		PhoneNumber: jobDetail.PhoneNumber,
		TrxId:       jobIdStr,
	}

	liveStatusResponse, err := svc.CreateLiveStatus(request, apiKey)
	if err != nil {
		svc.UpdateInvalidJobDetail(jobDetail.Id, "can not connected to server")

		return err
	}

	// todo: jika status code 200 kirim job detail ke aifcore
	// todo: jika status code 200 maka hapus job detail pada temp tabel. Sampai aifcore menyediakan API untuk get job details, untuk sementara jika status code 200 lakukan update subcriber_status dan device_status pada job detail
	if liveStatusResponse == nil {
		err = svc.UpdateFailedJobDetail(jobDetail.Id, jobDetail.Sequence)
		if err != nil {
			return err
		}
	} else {
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
			data.Scan(responseBodyByte)
		}

		if liveStatusResponse.StatusCode == 200 {
			// todo: pastikan errors bukan kode 6001, update kolom status "success", jika errors code 6001 update status "fail", hanya status "error" yg diulang
			// err = svc.DeleteJobDetail(jobDetail.Id)
			var status string
			if errorCode == -60001 {
				status = "fail"
			} else {
				status = "success"
			}
			err = svc.UpdateSucceededJobDetail(jobDetail.Id, subscriberStatus, deviceStatus, status, data)
			if err != nil {
				return err
			}
		} else {
			err = svc.UpdateFailedJobDetail(jobDetail.Id, jobDetail.Sequence)
			if err != nil {
				return err
			}
		}
	}

	return nil
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
	if liveStatusResponse != nil {
		liveStatusResponse.StatusCode = response.StatusCode
	}

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

func (svc *service) UpdateFailedJobDetail(id uint, sequence int) error {
	updateJobDetail := map[string]interface{}{}

	// maximumAttempts := 3
	// if sequence != maximumAttempts {
	// 	updateJobDetail["sequence"] = sequence + 1
	// 	updateJobDetail["on_process"] = true
	// 	updateJobDetail["status"] = "error"
	// } else {
	updateJobDetail["on_process"] = false
	// }

	return svc.Repo.UpdateJobDetail(id, updateJobDetail)
}

func (svc *service) UpdateInvalidJobDetail(id uint, errMessage string) error {
	updateJobDetail := map[string]interface{}{}

	updateJobDetail["status"] = "error"
	updateJobDetail["on_process"] = false
	updateJobDetail["message"] = errMessage

	return svc.Repo.UpdateJobDetail(id, updateJobDetail)
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

func (svc *service) GetJobDetailsByJobIdExport(jobId uint) ([]*JobDetailQueryResult, error) {
	jobDetails, err := svc.Repo.GetJobDetailsByJobIdExport(jobId)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (svc *service) GetJobWithIncompleteStatus() ([]uint, error) {
	return svc.Repo.GetJobWithIncompleteStatus()
}

func (svc *service) GetOnProcessJobDetails(jobId uint, onProcess bool) ([]uint, error) {
	return svc.Repo.GetOnProcessJobDetails(jobId, onProcess)
}

func (svc *service) CountOnProcessJobDetails(jobId uint, onProcess bool) (int, error) {
	count, err := svc.Repo.CountOnProcessJobDetails(jobId, onProcess)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
