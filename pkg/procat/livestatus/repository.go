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
	CreateJobInTx(userID string, dataJob *Job, dataJobDetail []LiveStatusRequest) (uint, error)
	GetJobs(limit, offset int, userID, startTime, endTime string) ([]*Job, error)
	GetJobsTotalByRangeDate(userID, startTime, endTime string) (int64, error)
	GetJobDetailsPercentageByDataAndRangeDate(userID, startTime, endTime, column, keyword string) (int64, error)
	GetJobByID(jobID uint) (*Job, error)
	GetJobByIDAndUserID(jobID uint, userID string) (*Job, error)
	GetJobsTotal(startTime, endTime string) (int64, error)
	GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error)
	GetJobDetailsByRangeDate(userID, startTime, endTime string) ([]*JobDetailQueryResult, error)
	GetJobDetailsByJobIDWithPagination(limit, offset int, keyword string, jobID uint) ([]*JobDetailQueryResult, error)
	GetJobDetailsByJobIDWithPaginationTotal(keyword string, jobID uint) (int64, error)
	GetJobDetailsByJobIDWithPaginationTotaPercentage(jobID uint, status string) (int64, error)
	GetJobDetailsTotalPercentageByStatusAndRangeDate(userID, startTime, endTime, status string) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error)
	GetFailedJobDetails() ([]*JobDetail, error)
	CallLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*http.Response, error)
	UpdateJob(id uint, req map[string]interface{}) error
	UpdateJobDetail(id uint, request map[string]interface{}) error
	DeleteJobDetail(id uint) error
	DeleteJob(id uint) error
}

func (repo *repository) CreateJobInTx(userID string, dataJob *Job, requests []LiveStatusRequest) (uint, error) {
	repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataJob).Error; err != nil {
			return err
		}

		for _, request := range requests {
			dataJobDetail := &JobDetail{
				UserID:      userID,
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

	return dataJob.ID, nil
}

func (repo *repository) GetJobs(limit, offset int, userID, startTime, endTime string) ([]*Job, error) {
	var jobs []*Job

	query := repo.DB.Where("user_id = ?", userID)
	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	if err := query.Limit(limit).Offset(offset).Order("id desc").Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (repo *repository) GetJobsTotalByRangeDate(userID, startTime, endTime string) (int64, error) {
	var totalData int64

	if err := repo.DB.Where("user_id = ? AND on_process = ? AND created_at BETWEEN ? AND ?", userID, false, startTime, endTime).Find(&JobDetail{}).Count(&totalData).Error; err != nil {
		return 0, err
	}

	return totalData, nil
}

func (repo *repository) GetJobDetailsPercentageByDataAndRangeDate(userID, startTime, endTime, column, keyword string) (int64, error) {
	var count int64

	query := repo.DB.Where("user_id = ? AND on_process = ? AND created_at BETWEEN ? AND ?", userID, false, startTime, endTime)

	if column == "subscriber_status" {
		query = query.Where("subscriber_status = ?", keyword)
	}

	if column == "device_status" {
		query = query.Where("device_status = ?", keyword)
	}

	if column == "data" && (keyword == "MOBILE" || keyword == "FIXED_LINE") {
		query = query.Where("data -> 'phone_type' ->> 'description' = ?", keyword)
	}

	err := query.Find(&JobDetail{}).Count(&count).Error
	return count, err
}

func (repo *repository) GetJobByID(jobID uint) (*Job, error) {
	var job *Job
	if err := repo.DB.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (repo *repository) GetJobByIDAndUserID(jobID uint, userID string) (*Job, error) {
	var job *Job
	if err := repo.DB.First(&job, "id = ? AND user_id = ?", jobID, userID).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (repo *repository) GetJobsTotal(startTime, endTime string) (int64, error) {
	var jobs []Job
	var count int64

	query := repo.DB
	if startTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	err := query.Find(&jobs).Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error) {
	var jobDetails []*JobDetail
	if err := repo.DB.Find(&jobDetails, "job_id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return jobDetails, nil
}

func (repo *repository) GetJobDetailsByRangeDate(userID, startTime, endTime string) ([]*JobDetailQueryResult, error) {
	var jobs []*JobDetailQueryResult
	err := repo.DB.
		Model(&JobDetail{}).
		Select("id, job_id, phone_number, subscriber_status, device_status, status, data -> 'carrier' ->> 'name' as operator, data -> 'phone_type' ->> 'description' as phone_type").
		Where("user_id = ? AND on_process = ? AND created_at BETWEEN ? AND ?", userID, false, startTime, endTime).
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

	query := repo.DB.Find(&jobs, "job_id = ? AND phone_number LIKE ?", jobID, "%"+keyword+"%")
	err := query.Find(&jobs).Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsByJobIDWithPaginationTotaPercentage(jobID uint, status string) (int64, error) {
	var jobs []JobDetail
	var count int64

	query := repo.DB.Find(&jobs, "job_id = ? and status = ?", jobID, status)
	err := query.Find(&jobs).Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsTotalPercentageByStatusAndRangeDate(userID, startTime, endTime, status string) (int64, error) {
	var count int64

	err := repo.DB.
		Where("user_id = ? AND on_process = ? AND created_at BETWEEN ? AND ? AND status = ?", userID, false, startTime, endTime, status).
		Find(&JobDetail{}).
		Count(&count).Error

	return count, err
}

func (repo *repository) GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error) {
	var jobs []JobDetail
	var count int64

	query := repo.DB.Where("job_id = ?", jobID)

	if column == "subscriber_status" {
		query = query.Where("subscriber_status = ?", keyword)
	}

	if column == "device_status" {
		query = query.Where("device_status = ?", keyword)
	}

	err := query.Find(&jobs).Count(&count).Error

	return count, err
}

func (repo *repository) GetFailedJobDetails() ([]*JobDetail, error) {
	var jobDetails []*JobDetail
	if err := repo.DB.Find(&jobDetails, "status = ?", "error").Error; err != nil {
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
