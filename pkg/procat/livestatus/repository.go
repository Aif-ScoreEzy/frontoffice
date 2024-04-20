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
	CreateJobInTx(dataJob *Job, dataJobDetail []LiveStatusRequest) (uint, error)
	GetJobs(limit, offset int) ([]*Job, error)
	GetJobByID(jobID uint) (*Job, error)
	GetJobsTotal() (int64, error)
	GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error)
	GetJobDetailsByJobIDWithPagination(limit, offset int, keyword string, jobID uint) ([]*JobDetail, error)
	GetJobDetailsByJobIDWithPaginationTotal(keyword string, jobID uint) (int64, error)
	GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error)
	GetUnprocessedJobDetails() ([]*JobDetail, error)
	CallLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*http.Response, error)
	UpdateJob(id uint, total int) error
	UpdateJobDetail(id uint, request *UpdateJobDetailRequest) error
	DeleteJobDetail(id uint) error
	DeleteJob(id uint) error
}

func (repo *repository) CreateJobInTx(dataJob *Job, requests []LiveStatusRequest) (uint, error) {
	repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataJob).Error; err != nil {
			return err
		}

		for _, request := range requests {
			dataJobDetail := &JobDetail{
				JobID:       dataJob.ID,
				PhoneNumber: request.PhoneNumber,
			}
			if err := tx.Create(dataJobDetail).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return dataJob.ID, nil
}

func (repo *repository) GetJobs(limit, offset int) ([]*Job, error) {
	var jobs []*Job
	if err := repo.DB.Limit(limit).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (repo *repository) GetJobByID(jobID uint) (*Job, error) {
	var job *Job
	if err := repo.DB.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (repo *repository) GetJobsTotal() (int64, error) {
	var jobs []Job
	var count int64

	query := repo.DB

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

func (repo *repository) GetJobDetailsByJobIDWithPagination(limit, offset int, keyword string, jobID uint) ([]*JobDetail, error) {
	var jobs []*JobDetail
	if err := repo.DB.Limit(limit).Offset(offset).Find(&jobs, "job_id = ? AND phone_number LIKE ?", jobID, "%"+keyword+"%").Error; err != nil {
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

func (repo *repository) GetJobDetailsPercentage(column, keyword string, jobID uint) (int64, error) {
	var jobs []JobDetail
	var count int64

	query := repo.DB.Find(&jobs, "job_id = ?", jobID)

	if column == "subscriber_status" {
		query = query.Where("subscriber_status = ?", keyword)
	}

	if column == "device_status" {
		query = query.Where("device_status = ?", keyword)
	}

	err := query.Find(&jobs).Count(&count).Error

	return count, err
}

func (repo *repository) GetUnprocessedJobDetails() ([]*JobDetail, error) {
	var jobDetails []*JobDetail
	if err := repo.DB.Find(&jobDetails, "on_process = ?", false).Error; err != nil {
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

func (repo *repository) UpdateJob(id uint, total int) error {
	if err := repo.DB.Model(&Job{}).Where("id = ?", id).Update("success", total).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repository) UpdateJobDetail(id uint, request *UpdateJobDetailRequest) error {
	if err := repo.DB.Model(&JobDetail{}).Where("id = ?", id).Updates(request).Error; err != nil {
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
