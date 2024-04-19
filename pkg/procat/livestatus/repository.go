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
	GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error)
	CallLiveStatus(liveStatusRequest *LiveStatusRequest, apiKey string) (*http.Response, error)
	UpdateJob(id uint, total int) error
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

func (repo *repository) GetJobDetailsByJobID(jobID uint) ([]*JobDetail, error) {
	var jobs []*JobDetail
	if err := repo.DB.Find(&jobs, "job_id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return jobs, nil
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
