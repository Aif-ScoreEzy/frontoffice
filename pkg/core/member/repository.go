package member

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{
		Db:  db,
		Cfg: cfg,
	}
}

type repository struct {
	Db  *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	GetMemberBy(query *FindUserQuery) (*http.Response, error)
	GetMemberList() (*http.Response, error)
	DeleteMemberById(id string) (*http.Response, error)
}

func (repo *repository) GetMemberBy(query *FindUserQuery) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/by`, repo.Cfg.Env.AifcoreHost)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("id", query.Id)
	q.Add("email", query.Email)
	q.Add("username", query.Username)
	q.Add("key", query.Key)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) GetMemberList() (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/list`, repo.Cfg.Env.AifcoreHost)
	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) DeleteMemberById(id string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/deletemember/%v`, repo.Cfg.Env.AifcoreHost, id)
	request, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	return client.Do(request)
}
