package genretail

import (
	"front-office/config/database"

	"gorm.io/gorm"
)

func StoreImportData(newData []*BulkSearch, userID string) error {
	errTx := database.DBConn.Transaction(func(tx *gorm.DB) error {
		// remove data existing in table
		if err := database.DBConn.Debug().Delete(&BulkSearch{}, "user_id = ?", userID).Error; err != nil {
			return err
		}

		// replace existing with new data
		if err := tx.Debug().Create(&newData).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return errTx
	}

	return nil
}

func GetAllBulkSearch(tierLevel uint, userID, companyID string) ([]BulkSearch, error) {
	var bulkSearches []BulkSearch

	query := database.DBConn.Debug()

	if tierLevel == 1 {
		// admin
		query = query.Where("company_id = ?", companyID)
	} else {
		// user
		query = query.Where("user_id = ?", userID)
	}

	err := query.Find(&bulkSearches)

	if err.Error != nil {
		return nil, err.Error
	}

	return bulkSearches, nil
}

func CountData(tierLevel uint, userID, companyID string) (int64, error) {
	var bulkSearches []BulkSearch
	var count int64

	query := database.DBConn.Debug()

	if tierLevel == 1 {
		// admin
		query = query.Where("company_id = ?", companyID)
	} else {
		// user
		query = query.Where("user_id = ?", userID)
	}

	err := query.Find(&bulkSearches).Count(&count).Error

	return count, err
}
