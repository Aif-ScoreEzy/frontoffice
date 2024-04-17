package livestatus

import (
	"front-office/app/config"

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
