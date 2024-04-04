package livestatus

import (
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

type repository struct {
	DB *gorm.DB
}

type Repository interface {
	CreateJobInTx(dataJob *Job, dataJobDetail *FIFRequests) error
}

func (repo *repository) CreateJobInTx(dataJob *Job, requests *FIFRequests) error {
	repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataJob).Error; err != nil {
			return err
		}

		for _, request := range requests.PhoneNumbers {
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

	return nil
}
