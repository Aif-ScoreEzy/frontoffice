package member

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"mime/multipart"
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
	AddMember(req *RegisterMemberRequest) (*http.Response, error)
	GetMemberBy(query *FindUserQuery) (*http.Response, error)
	GetMemberList(filter *MemberFilter) (*http.Response, error)
	UpdateOneById(id string, req map[string]interface{}) (*http.Response, error)
	DeleteMemberById(id string) (*http.Response, error)
}

func (repo *repository) AddMember(req *RegisterMemberRequest) (*http.Response, error) {
	apiUrl := repo.Cfg.Env.AifcoreHost + "/api/core/member/addmember"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("name", req.Name)
	writer.WriteField("email", req.Email)
	writer.WriteField("key", req.Key)
	writer.WriteField("companyid", fmt.Sprintf("%d", req.CompanyId))

	writer.Close()

	request, _ := http.NewRequest(http.MethodPost, apiUrl, &body)
	request.Header.Set(constant.HeaderContentType, writer.FormDataContentType())
	request.Header.Set(constant.XAPIKey, repo.Cfg.Env.XModuleKey)

	client := &http.Client{}
	return client.Do(request)
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

func (repo *repository) GetMemberList(filter *MemberFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/listbycompany/%v`, repo.Cfg.Env.AifcoreHost, filter.CompanyID)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("page", filter.Page)
	q.Add("size", filter.Limit)
	q.Add("keyword", filter.Keyword)
	q.Add("status", filter.Status)
	q.Add("role_id", filter.RoleID)
	q.Add("start_date", filter.StartDate)
	q.Add("end_date", filter.EndDate)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) UpdateOneById(id string, req map[string]interface{}) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/updateprofile/%v`, repo.Cfg.Env.AifcoreHost, id)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(jsonBodyValue))
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
