package grading

import "github.com/google/uuid"

func CreateGradingsSvc(req *CreateGradingRequest, companyID string) (*Grading, error) {
	gradingCompanyID := uuid.NewString()

	gradingData := &Grading{
		ID:           gradingCompanyID,
		GradingLabel: req.GradingLabel,
		MinGrade:     req.MinGrade,
		MaxGrade:     req.MaxGrade,
		CompanyID:    companyID,
	}

	grading, err := CreateGrading(gradingData)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func GetGradingByGradinglabelSvc(gradingLabel, companyID string) (*Grading, error) {
	grading, err := FindOneByGradingLabel(gradingLabel, companyID)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func GetGradingsSvc(companyID string) ([]*Grading, error) {
	gradings, err := FindAllGradings(companyID)
	if err != nil {
		return nil, err
	}

	return gradings, nil
}
