package grading

import "front-office/config/database"

func CreateGrading(grading *Grading) (*Grading, error) {
	query := database.DBConn.Debug().Create(&grading)
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
