package livestatus

import (
	"bytes"
	"encoding/json"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{DB: db, Cfg: cfg}
}

type repository struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	CreateJobInTx(userID, companyID string, dataJob *Job, dataJobDetail []LiveStatusRequest) (uint, error)
	GetJobs(limit, offset int, tierLevel uint, userID, companyID, startTime, endTime string) ([]*Job, error)
	GetJobsTotalByRangeDate(userID, companyID, startTime, endTime string, tierLevel uint) (int64, error)
	GetJobDetailsPercentageByDataAndRangeDate(userID, companyID, startTime, endTime, column, keyword string, tierLevel uint) (int64, error)
	GetJobByID(jobID uint) (*Job, error)
	GetJobByIDAndUserID(jobID, tierLevel uint, userID, companyID string) (*Job, error)
	GetJobsTotal(userID, companyID, startTime, endTime string, tierLevel uint) (int64, error)
	GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error)
	GetJobDetailsByRangeDate(userID, companyID, startTime, endTime string, tierLevel uint) ([]*JobDetailQueryResult, error)
	GetJobDetailsByJobIDWithPagination(limit, offset int, keyword string, jobID uint) ([]*JobDetailQueryResult, error)
	GetJobDetailsByJobIDWithPaginationTotal(keyword string, jobID uint) (int64, error)
	GetJobDetailsByJobIDWithPaginationTotaPercentage(jobID uint, status string) (int64, error)
	GetJobDetailsTotalPercentageByStatusAndRangeDate(userID, companyID, startTime, endTime, status string, tierLevel uint) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error)
	GetFailedJobDetails(jobID uint) ([]*JobDetail, error)
	CallLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*http.Response, error)
	UpdateJob(id uint, req map[string]interface{}) error
	UpdateJobDetail(id uint, request map[string]interface{}) error
	DeleteJobDetail(id uint) error
	DeleteJob(id uint) error
	GetJobDetailsByJobIDExport(jobID uint) ([]*JobDetailQueryResult, error)
	GetJobWithIncompleteStatus() ([]uint, error)
	GetOnProcessJobDetails(jobID uint, onProcess bool) ([]uint, error)
	CountOnProcessJobDetails(jobID uint, onProcess bool) (int64, error)
}

func (repo *repository) CreateJobInTx(userId, companyId string, dataJob *Job, requests []LiveStatusRequest) (uint, error) {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataJob).Error; err != nil {
			return err
		}

		for _, request := range requests {
			dataJobDetail := &JobDetail{
				UserID:      userID,
				CompanyID:   companyID,
				JobID:       dataJob.ID,
				PhoneNumber: request.PhoneNumber,
				OnProcess:   true,
			}
			if err := tx.Create(dataJobDetail).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if errTx != nil {
		return 0, errTx
	}

	return dataJob.Id, nil
}

func (repo *repository) GetJobs(limit, offset int, tierLevel uint, userID, companyID, startTime, endTime string) ([]*Job, error) {
	var jobs []*Job

	query := repo.DB

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	if err := query.Limit(limit).Offset(offset).Order("id desc").Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (repo *repository) GetJobsTotalByRangeDate(userID, companyID, startTime, endTime string, tierLevel uint) (int64, error) {
	var totalData int64

	query := repo.DB.Model(&JobDetail{})

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Where("on_process = ? AND created_at BETWEEN ? AND ?", false, startTime, endTime).Count(&totalData).Error; err != nil {
		return 0, err
	}

	return totalData, nil
}

func (repo *repository) GetJobDetailsPercentageByDataAndRangeDate(userID, companyID, startTime, endTime, column, keyword string, tierLevel uint) (int64, error) {
	var count int64

	query := repo.DB.Model(&JobDetail{})

	query = query.Where("on_process = ? AND created_at BETWEEN ? AND ?", false, startTime, endTime)

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if column == "subscriber_status" {
		query = query.Where("subscriber_status = ?", keyword)
	}

	if column == "device_status" {
		query = query.Where("device_status = ?", keyword)
	}

	if column == "data" && (keyword == "MOBILE" || keyword == "FIXED_LINE") {
		query = query.Where("data -> 'phone_type' ->> 'description' = ?", keyword)
	}

	err := query.Count(&count).Error
	return count, err
}

func (repo *repository) GetJobByID(jobID uint) (*Job, error) {
	var job *Job
	if err := repo.DB.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (repo *repository) GetJobByIDAndUserID(jobID, tierLevel uint, userID, companyID string) (*Job, error) {
	var job *Job

	query := repo.DB

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (repo *repository) GetJobsTotal(userID, companyID, startTime, endTime string, tierLevel uint) (int64, error) {
	var jobs []Job
	var count int64

	query := repo.DB.Model(&jobs)

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	err := query.Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error) {
	var jobDetails []*JobDetail
	if err := repo.DB.Find(&jobDetails, "job_id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (repo *repository) GetJobDetailsByRangeDate(userID, companyID, startTime, endTime string, tierLevel uint) ([]*JobDetailQueryResult, error) {
	var jobs []*JobDetailQueryResult

	query := repo.DB

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	err := query.
		Model(&JobDetail{}).
		Select("id, job_id, phone_number, subscriber_status, device_status, status, data -> 'carrier' ->> 'name' as operator, data -> 'phone_type' ->> 'description' as phone_type, message").
		Where("on_process = ? AND created_at BETWEEN ? AND ?", false, startTime, endTime).
		Find(&jobs).
		Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (repo *repository) GetJobDetailsByJobIDWithPagination(limit, offset int, keyword string, jobID uint) ([]*JobDetailQueryResult, error) {
	var jobs []*JobDetailQueryResult

	if err := repo.DB.
		Model(&JobDetail{}).
		Select("id, job_id, phone_number, subscriber_status, device_status, status, data -> 'carrier' ->> 'name' as operator, data -> 'phone_type' ->> 'description' as phone_type").
		Limit(limit).
		Offset(offset).
		Where("job_id = ? AND phone_number LIKE ?", jobID, "%"+keyword+"%").
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (repo *repository) GetJobDetailsByJobIDWithPaginationTotal(keyword string, jobID uint) (int64, error) {
	var jobs []JobDetail
	var count int64

	err := repo.DB.Model(&jobs).Where("job_id = ? AND phone_number LIKE ?", jobID, "%"+keyword+"%").Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsByJobIDWithPaginationTotaPercentage(jobID uint, status string) (int64, error) {
	var jobs []JobDetail
	var count int64

	query := repo.DB.Model(&jobs).Where("job_id = ? and status = ?", jobID, status)
	err := query.Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsTotalPercentageByStatusAndRangeDate(userID, companyID, startTime, endTime, status string, tierLevel uint) (int64, error) {
	var count int64

	query := repo.DB.Model(&JobDetail{})

	if tierLevel == 1 {
		query = query.Where("company_id = ?", companyID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	err := query.
		Where("on_process = ? AND created_at BETWEEN ? AND ? AND status = ?", false, startTime, endTime, status).
		Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error) {
	var jobs []JobDetail
	var count int64

	query := repo.DB.Model(&jobs).Where("job_id = ?", jobID)

	if column == "subscriber_status" {
		query = query.Where("subscriber_status = ?", keyword)
	}

	if column == "device_status" {
		query = query.Where("device_status = ?", keyword)
	}

	err := query.Count(&count).Error

	return count, err
}

func (repo *repository) GetFailedJobDetails(jobID uint) ([]*JobDetail, error) {
	var jobDetails []*JobDetail
	maximumAttempts := 3
	if err := repo.DB.Find(&jobDetails, "job_id = ? AND status = ? AND sequence <= ? AND on_process = true", jobID, "error", maximumAttempts).Error; err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (repo *repository) CallLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.PartnerServiceHost + "/api/partner/telesign/phone-live-status"

	jsonBodyValue, _ := json.Marshal(liveStatusRequest)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	request.Header.Set("X-AIF-KEY", apiKey)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) UpdateJob(id uint, req map[string]interface{}) error {
	if err := repo.DB.Model(&Job{}).Where("id = ?", id).Updates(req).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repository) UpdateJobDetail(id uint, data map[string]interface{}) error {
	if err := repo.DB.Model(&JobDetail{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repository) DeleteJobDetail(id uint) error {
	if err := repo.DB.Delete(&JobDetail{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repository) DeleteJob(id uint) error {
	if err := repo.DB.Delete(&Job{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetJobDetailsByJobIDExport(jobID uint) ([]*JobDetailQueryResult, error) {
	var jobs []*JobDetailQueryResult

	if err := repo.DB.
		Model(&JobDetail{}).
		Select("id, job_id, phone_number, subscriber_status, device_status, status, data -> 'carrier' ->> 'name' as operator, data -> 'phone_type' ->> 'description' as phone_type").
		Where("job_id = ?", jobID).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (repo *repository) GetJobWithIncompleteStatus() ([]uint, error) {
	var jobIDs []uint
	if err := repo.DB.Model(&Job{}).Select("id").Where("status = ?", "").Find(&jobIDs).Error; err != nil {
		return nil, err
	}

	return jobIDs, nil
}

func (repo *repository) GetOnProcessJobDetails(jobID uint, onProcess bool) ([]uint, error) {
	var jobDetailIDs []uint
	if err := repo.DB.Model(&JobDetail{}).Select("id").Where("job_id = ? AND on_process = ?", jobID, onProcess).Find(&jobDetailIDs).Error; err != nil {
		return nil, err
	}

	return jobDetailIDs, nil
}

func (repo *repository) CountOnProcessJobDetails(jobID uint, onProcess bool) (int64, error) {
	var count int64

	if err := repo.DB.Model(&JobDetail{}).Where("job_id = ? AND on_process = ?", jobID, onProcess).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
