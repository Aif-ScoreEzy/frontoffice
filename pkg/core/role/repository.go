package role

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
)

func NewRepository(cfg *config.Config) Repository {
	return &repository{
		Cfg: cfg,
	}
}

type repository struct {
	Cfg *config.Config
}

type Repository interface {
	FindAll(filter RoleFilter) (*http.Response, error)
	FindOneById(id string) (*http.Response, error)
}

func (repo *repository) FindOneById(id string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/role/%v`, repo.Cfg.Env.AifcoreHost, id)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) FindAll(filter RoleFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/role`, repo.Cfg.Env.AifcoreHost)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("name", filter.Name)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}
