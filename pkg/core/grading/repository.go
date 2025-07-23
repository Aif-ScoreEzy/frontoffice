package grading

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
)

func NewRepository(cfg *config.Config) Repository {
	return &repository{
		cfg: cfg,
	}
}

type repository struct {
	cfg *config.Config
}

type Repository interface {
	GetGradeList(apiconfigId string) (*http.Response, error)
}

func (repo *repository) GetGradeList(apiconfigId string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/apiconfig`, repo.cfg.Env.AifcoreHost)

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
