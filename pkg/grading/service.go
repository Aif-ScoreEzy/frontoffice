package grading

import (
	"errors"
	"front-office/constant"
	"time"

	"github.com/google/uuid"
)

func CreateGradingSvc(req *CreateGradingRequest, companyID string) (*Grading, error) {
	gradingID := uuid.NewString()

	if req.GradingLabel == "" {
		return nil, errors.New(constant.FieldGradingLabelEmpty)
	}

	if req.MinGrade == nil {
		return nil, errors.New(constant.FieldMinGradeEmpty)
	}

	if req.MaxGrade == nil {
		return nil, errors.New(constant.FieldMaxGradeEmpty)
	}

	gradingData := &Grading{
		ID:           gradingID,
		GradingLabel: req.GradingLabel,
		MinGrade:     *req.MinGrade,
		MaxGrade:     *req.MaxGrade,
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

func GetGradingByIDSvc(gradingID, companyID string) (*Grading, error) {
	grading, err := FindOneByID(gradingID, companyID)
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

func UpdateGradingSvc(req *UpdateGradingRequest, companyID string) (*Grading, error) {
	updateGrading := &UpdateGradingRequest{}

	if req.IsDeleted {
		updateGrading = &UpdateGradingRequest{
			DeletedAt: time.Now(),
		}
	} else {
		updateGrading = &UpdateGradingRequest{
			GradingLabel: req.GradingLabel,
			MinGrade:     req.MinGrade,
			MaxGrade:     req.MaxGrade,
			UpdatedAt:    time.Now(),
		}
	}

	grading, err := UpdateOneByID(updateGrading, req.ID, companyID)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func DeleteGradingsSvc(companyID string) error {
	err := DeleteAllGradings(companyID)
	if err != nil {
		return err
	}

	return nil
}
