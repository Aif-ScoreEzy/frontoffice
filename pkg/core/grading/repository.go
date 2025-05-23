package grading

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{
		DB:  db,
		Cfg: cfg,
	}
}

type repository struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	GetGradeList(apiconfigId string) (*http.Response, error)
	CreateGrading(grading *Grading) (*Grading, error)
	FindOneById(gradingId, companyId string) (*Grading, error)
	FindOneByGradingLabel(gradingLabel, companyId string) (*Grading, error)
	FindAllGradings(companyId string) ([]*Grading, error)
	UpdateOneById(updateGrading *UpdateGradingRequest, gradingId, companyId string) (*Grading, error)
	ReplaceAllGradings(gradings []*Grading, companyId string) error
	DeleteAllGradings(companyId string) error
}

func (repo *repository) GetGradeList(apiconfigId string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/apiconfig`, repo.Cfg.Env.AifcoreHost)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("id", apiconfigId)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) CreateGrading(grading *Grading) (*Grading, error) {
	query := repo.DB.Create(&grading)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) FindOneById(gradingId, companyId string) (*Grading, error) {
	var grading *Grading

	query := repo.DB.First(&grading, "id = ? AND company_id = ?", gradingId, companyId)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) FindOneByGradingLabel(gradingLabel, companyId string) (*Grading, error) {
	var grading *Grading

	query := repo.DB.First(&grading, "grading_label = ? AND company_id = ?", gradingLabel, companyId)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) FindAllGradings(companyId string) ([]*Grading, error) {
	var gradings []*Grading

	query := repo.DB.Find(&gradings, "company_id = ?", companyId)
	if query.Error != nil {
		return nil, query.Error
	}

	return gradings, nil
}

func (repo *repository) UpdateOneById(updateGrading *UpdateGradingRequest, gradingId, companyId string) (*Grading, error) {
	var grading *Grading

	query := repo.DB.Model(&grading).Where("id = ? AND company_id = ?", gradingId, companyId).Updates(updateGrading)
	if query.Error != nil {
		return nil, query.Error
	}

	return grading, nil
}

func (repo *repository) ReplaceAllGradings(gradings []*Grading, companyId string) error {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := repo.DB.Delete(&Grading{}, "company_id = ?", companyId).Error; err != nil {
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

func (repo *repository) DeleteAllGradings(companyId string) error {
	query := repo.DB.Delete(&Grading{}, "company_id = ?", companyId)
	if query.Error != nil {
		return query.Error
	}

	return nil
}
