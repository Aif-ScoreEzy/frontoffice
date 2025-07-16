package template

import (
	"front-office/common/constant"
)

type Service interface {
	ListTemplates() ([]TemplateInfo, error)
	DownloadTemplate(req DownloadRequest) (string, error)
}

type service struct {
	Repo Repository
}

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}
func (s *service) DownloadTemplate(req DownloadRequest) (string, error) {
	return s.Repo.GetTemplatePath(req.Category, req.Filename)
}

func (s *service) ListTemplates() ([]TemplateInfo, error) {
	templates, err := s.Repo.GetAvailableTemplates()
	if err != nil {
		return nil, err
	}

	var result []TemplateInfo
	for category, files := range templates {
		desc := getCategoryDescription(category)
		result = append(result, TemplateInfo{
			Category:    category,
			Files:       files,
			Description: desc,
		})
	}

	return result, nil
}

func getCategoryDescription(category string) string {
	switch category {
	case constant.PhoneLiveTemplates:
		return "Phone live status template"
	case constant.LoanRecordCheckerTemplates:
		return "Loan record checker template"
	case constant.MultipleLoanTemplates:
		return "Multiple loan template"
	case constant.TaxComplianceStatusTemplates:
		return "Tax compliance status template"
	case constant.TaxScoreTemplates:
		return "Tax score template"
	default:
		return "Common template"
	}
}
