package genretail

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
	StoreImportData(newData []*BulkSearch, userId string) error
	GetAllBulkSearch(tierLevel uint, userId, companyId string) ([]*BulkSearch, error)
	CountData(tierLevel uint, userId, companyId string) (int64, error)
}

func (repo *repository) StoreImportData(newData []*BulkSearch, userId string) error {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		// remove data existing in table
		if err := repo.DB.Delete(&BulkSearch{}, "user_id = ?", userId).Error; err != nil {
			return err
		}

		// replace existing with new data
		if err := tx.Create(&newData).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return errTx
	}

	return nil
}

func (repo *repository) GetAllBulkSearch(tierLevel uint, userId, companyId string) ([]*BulkSearch, error) {
	var bulkSearches []*BulkSearch

	query := repo.DB.Preload("User")

	if tierLevel == 1 {
		// admin
		query = query.Where("company_id = ?", companyId)
	} else {
		// user
		query = query.Where("user_id = ?", userId)
	}

	err := query.Find(&bulkSearches)

	if err.Error != nil {
		return nil, err.Error
	}

	return bulkSearches, nil
}

func (repo *repository) CountData(tierLevel uint, userId, companyId string) (int64, error) {
	var bulkSearches []*BulkSearch
	var count int64

	query := repo.DB.Debug()

	if tierLevel == 1 {
		// admin
		query = query.Where("company_id = ?", companyId)
	} else {
		// user
		query = query.Where("user_id = ?", userId)
	}

	err := query.Find(&bulkSearches).Count(&count).Error

	return count, err
}
