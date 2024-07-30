package grading

import (
	"errors"
	"front-office/common/constant"
	"time"

	"github.com/google/uuid"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateGradingSvc(req *CreateGradingRequest, companyId string) (*Grading, error)
	GetGradingByGradinglabelSvc(gradingLabel, companyId string) (*Grading, error)
	GetGradingByIdSvc(gradingId, companyId string) (*Grading, error)
	GetGradingsSvc(companyId string) ([]*Grading, error)
	UpdateGradingSvc(req *UpdateGradingRequest, companyId string) (*Grading, error)
	ReplaceAllGradingsSvc(createGradingsRequest *CreateGradingsRequest, companyId string) error
	ReplaceAllGradingsNewSvc(createGradingsRequest *CreateGradingsNewRequest, companyId string) error
	DeleteGradingsSvc(companyId string) error
}

func (svc *service) CreateGradingSvc(req *CreateGradingRequest, companyId string) (*Grading, error) {
	gradingId := uuid.NewString()

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
		Id:           gradingId,
		GradingLabel: req.GradingLabel,
		MinGrade:     *req.MinGrade,
		MaxGrade:     *req.MaxGrade,
		CompanyId:    companyId,
	}

	grading, err := svc.Repo.CreateGrading(gradingData)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func (svc *service) GetGradingByGradinglabelSvc(gradingLabel, companyId string) (*Grading, error) {
	grading, err := svc.Repo.FindOneByGradingLabel(gradingLabel, companyId)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func (svc *service) GetGradingByIdSvc(gradingId, companyId string) (*Grading, error) {
	grading, err := svc.Repo.FindOneById(gradingId, companyId)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func (svc *service) GetGradingsSvc(companyId string) ([]*Grading, error) {
	gradings, err := svc.Repo.FindAllGradings(companyId)
	if err != nil {
		return nil, err
	}

	return gradings, nil
}

func (svc *service) UpdateGradingSvc(req *UpdateGradingRequest, companyId string) (*Grading, error) {
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

	grading, err := svc.Repo.UpdateOneById(updateGrading, req.Id, companyId)
	if err != nil {
		return nil, err
	}

	return grading, nil
}

func (svc *service) ReplaceAllGradingsSvc(createGradingsRequest *CreateGradingsRequest, companyId string) error {
	var gradings []*Grading

	for _, createGradingRequest := range createGradingsRequest.CreateGradingsRequest {
		if createGradingRequest.GradingLabel == "" {
			return errors.New(constant.FieldGradingLabelEmpty)
		}

		if createGradingRequest.MinGrade == nil {
			return errors.New(constant.FieldMinGradeEmpty)
		}

		if createGradingRequest.MaxGrade == nil {
			return errors.New(constant.FieldMaxGradeEmpty)
		}

		gradingId := uuid.NewString()
		grading := &Grading{
			Id:           gradingId,
			GradingLabel: createGradingRequest.GradingLabel,
			MinGrade:     *createGradingRequest.MinGrade,
			MaxGrade:     *createGradingRequest.MaxGrade,
			CompanyId:    companyId,
		}

		gradings = append(gradings, grading)
	}

	err := svc.Repo.ReplaceAllGradings(gradings, companyId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) ReplaceAllGradingsNewSvc(createGradingsRequest *CreateGradingsNewRequest, companyId string) error {
	var gradings []*Grading

	for _, createGradingRequest := range createGradingsRequest.CreateGradingsNewRequest {
		if createGradingRequest.Grade == "" {
			return errors.New(constant.FieldGradingLabelEmpty)
		}

		if len(createGradingRequest.Value) == 0 {
			return errors.New(constant.FieldGradingValueEmpty)
		}

		gradingId := uuid.NewString()
		// create the grading to append to the gradings
		grading := &Grading{
			Id:           gradingId,
			GradingLabel: createGradingRequest.Grade,
		}
		for i, v := range createGradingRequest.Value {
			if i == 0 {
				grading.MinGrade = *v
			}
			if i == 1 {
				grading.MaxGrade = *v
			}
		}
		grading.CompanyId = companyId
		gradings = append(gradings, grading)
	}

	err := svc.Repo.ReplaceAllGradings(gradings, companyId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) DeleteGradingsSvc(companyId string) error {
	err := svc.Repo.DeleteAllGradings(companyId)
	if err != nil {
		return err
	}

	return nil
}
