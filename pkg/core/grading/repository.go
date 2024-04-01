package grading

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
	CreateGrading(grading *Grading) (*Grading, error)
	FindOneByID(gradingID, companyID string) (*Grading, error)
	FindOneByGradingLabel(gradingLabel, companyID string) (*Grading, error)
	FindAllGradings(companyID string) ([]*Grading, error)
	UpdateOneByID(updateGrading *UpdateGradingRequest, gradingID, companyID string) (*Grading, error)
	ReplaceAllGradings(gradings []*Grading, companyID string) error
	DeleteAllGradings(companyID string) error
}

func (repo *repository) CreateGrading(grading *Grading) (*Grading, error) {
	query := repo.DB.Create(&grading)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) FindOneByID(gradingID, companyID string) (*Grading, error) {
	var grading *Grading

	query := repo.DB.First(&grading, "id = ? AND company_id = ?", gradingID, companyID)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) FindOneByGradingLabel(gradingLabel, companyID string) (*Grading, error) {
	var grading *Grading

	query := repo.DB.First(&grading, "grading_label = ? AND company_id = ?", gradingLabel, companyID)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) FindAllGradings(companyID string) ([]*Grading, error) {
	var gradings []*Grading

	query := repo.DB.Find(&gradings, "company_id = ?", companyID)
	if query.Error != nil {
		return nil, query.Error
	}

	return gradings, nil
}

func (repo *repository) UpdateOneByID(updateGrading *UpdateGradingRequest, gradingID, companyID string) (*Grading, error) {
	var grading *Grading

	query := repo.DB.Model(&grading).Where("id = ? AND company_id = ?", gradingID, companyID).Updates(updateGrading)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) ReplaceAllGradings(gradings []*Grading, companyID string) error {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := repo.DB.Delete(&Grading{}, "company_id = ?", companyID).Error; err != nil {
			return err
		}

		if err := tx.Create(&gradings).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return errTx
	}

	return nil
}

func (repo *repository) DeleteAllGradings(companyID string) error {
	query := repo.DB.Delete(&Grading{}, "company_id = ?", companyID)
	if query.Error != nil {
		return query.Error
	}

	return nil
}
