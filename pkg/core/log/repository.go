package log

import (
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{
		DB: db,
		Cfg: cfg,
	}
}

type repository struct {
	DB *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	FindAllTransactionLogsByDate(companyId, date string) (*http.Response, error)
}

func (repo *repository) FindAllTransactionLogsByDate(companyId, date string) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/by"

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("date", date)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}