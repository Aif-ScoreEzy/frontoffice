package grading

import (
	"front-office/config/database"

	"gorm.io/gorm"
)

func CreateGrading(grading *Grading) (*Grading, error) {
	query := database.DBConn.Debug().Create(&grading)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func FindOneByID(gradingID, companyID string) (*Grading, error) {
	var grading *Grading

	query := database.DBConn.Debug().First(&grading, "id = ? AND company_id = ?", gradingID, companyID)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func FindOneByGradingLabel(gradingLabel, companyID string) (*Grading, error) {
	var grading *Grading

	query := database.DBConn.Debug().First(&grading, "grading_label = ? AND company_id = ?", gradingLabel, companyID)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func FindAllGradings(companyID string) ([]*Grading, error) {
	var gradings []*Grading

	query := database.DBConn.Debug().Find(&gradings, "company_id = ?", companyID)
	if query.Error != nil {
		return nil, query.Error
	}

	return gradings, nil
}

func UpdateOneByID(updateGrading *UpdateGradingRequest, gradingID, companyID string) (*Grading, error) {
	var grading *Grading

	query := database.DBConn.Debug().Model(&grading).Where("id = ? AND company_id = ?", gradingID, companyID).Updates(updateGrading)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func ReplaceAllGradings(gradings []*Grading, companyID string) error {
	errTx := database.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := database.DBConn.Debug().Delete(&Grading{}, "company_id = ?", companyID).Error; err != nil {
			return err
		}

		if err := tx.Debug().Create(&gradings).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return errTx
	}

	return nil
}

func DeleteAllGradings(companyID string) error {
	query := database.DBConn.Debug().Delete(&Grading{}, "company_id = ?", companyID)
	if query.Error != nil {
		return query.Error
	}

	return nil
}
