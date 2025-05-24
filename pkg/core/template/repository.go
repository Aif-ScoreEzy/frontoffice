package template

import (
	"front-office/common/constant"
	"os"
	"path/filepath"
)

func NewRepository() Repository {
	return &repository{}
}

type repository struct{}

type Repository interface {
	GetAvailableTemplates() (map[string][]string, error)
	GetTemplatePath(category, filename string) (string, error)
}

func (r *repository) GetTemplatePath(category, filename string) (string, error) {
	path := filepath.Join(constant.TemplateBaseDir, category, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}
	return path, nil
}

func (r *repository) GetAvailableTemplates() (map[string][]string, error) {
	categories := []string{
		constant.PhoneLiveTemplates,
	}

	result := make(map[string][]string)

	for _, category := range categories {
		files, err := os.ReadDir(filepath.Join(constant.TemplateBaseDir, category))
		if err != nil {
			return nil, err
		}

		var filenames []string
		for _, file := range files {
			filenames = append(filenames, file.Name())
		}

		result[category] = filenames
	}

	return result, nil
}
