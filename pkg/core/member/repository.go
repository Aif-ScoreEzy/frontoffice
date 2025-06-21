package member

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/httpclient"
	"mime/multipart"
	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient) Repository {
	return &repository{cfg, client}
}

type repository struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type Repository interface {
	CallAddMemberAPI(req *RegisterMemberRequest) (*registerResponseData, error)
	GetMemberBy(query *FindUserQuery) (*http.Response, error)
	GetMemberList(filter *MemberFilter) (*http.Response, error)
	CallUpdateMemberAPI(id string, req map[string]interface{}) error
	DeleteMemberById(id string) (*http.Response, error)
}

func (repo *repository) CallAddMemberAPI(reqBody *RegisterMemberRequest) (*registerResponseData, error) {
	url := fmt.Sprintf("%s/api/core/member/addmember", repo.cfg.Env.AifcoreHost)

	var bodyBytes bytes.Buffer
	writer := multipart.NewWriter(&bodyBytes)

	writer.WriteField("name", reqBody.Name)
	writer.WriteField("email", reqBody.Email)
	writer.WriteField("key", reqBody.Key)
	writer.WriteField("companyid", fmt.Sprintf("%d", reqBody.CompanyId))
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, &bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set(constant.HeaderContentType, writer.FormDataContentType())
	req.Header.Set(constant.XAPIKey, repo.cfg.Env.XModuleKey)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*registerResponseData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) GetMemberBy(query *FindUserQuery) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/by`, repo.cfg.Env.AifcoreHost)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("id", query.Id)
	q.Add("company_id", query.CompanyId)
	q.Add("email", query.Email)
	q.Add("username", query.Username)
	q.Add("key", query.Key)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) GetMemberList(filter *MemberFilter) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/listbycompany/%v`, repo.cfg.Env.AifcoreHost, filter.CompanyID)

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

func (repo *repository) CallUpdateMemberAPI(id string, reqBody map[string]interface{}) error {
	url := fmt.Sprintf(`%v/api/core/member/updateprofile/%v`, repo.cfg.Env.AifcoreHost, id)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) DeleteMemberById(id string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/deletemember/%v`, repo.cfg.Env.AifcoreHost, id)
	request, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	return client.Do(request)
}
