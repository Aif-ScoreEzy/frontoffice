package member

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/httpclient"
	"front-office/internal/jsonutil"
	"mime/multipart"
	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient, marshalFn jsonutil.Marshaller) Repository {
	if marshalFn == nil {
		marshalFn = json.Marshal
	}

	return &repository{
		cfg:       cfg,
		client:    client,
		marshalFn: marshalFn,
	}
}

type repository struct {
	cfg       *config.Config
	client    httpclient.HTTPClient
	marshalFn jsonutil.Marshaller
}

type Repository interface {
	AddMemberAPI(req *RegisterMemberRequest) (*registerResponseData, error)
	GetMemberAPI(query *FindUserQuery) (*MstMember, error)
	GetMemberListAPI(filter *MemberFilter) ([]*MstMember, *model.Meta, error)
	CallUpdateMemberAPI(id string, req map[string]interface{}) error
	CallDeleteMemberAPI(id string) error
}

func (repo *repository) AddMemberAPI(reqBody *RegisterMemberRequest) (*registerResponseData, error) {
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
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	req.Header.Set(constant.HeaderContentType, writer.FormDataContentType())
	req.Header.Set(constant.XAPIKey, repo.cfg.Env.XModuleKey)

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*registerResponseData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) GetMemberAPI(query *FindUserQuery) (*MstMember, error) {
	url := fmt.Sprintf(`%v/api/core/member/by`, repo.cfg.Env.AifcoreHost)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := req.URL.Query()
	q.Add("id", query.Id)
	q.Add("company_id", query.CompanyId)
	q.Add("email", query.Email)
	q.Add("username", query.Username)
	q.Add("key", query.Key)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[*MstMember](resp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (repo *repository) GetMemberListAPI(filter *MemberFilter) ([]*MstMember, *model.Meta, error) {
	url := fmt.Sprintf(`%v/api/core/member/listbycompany/%v`, repo.cfg.Env.AifcoreHost, filter.CompanyID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := req.URL.Query()
	q.Add(constant.Page, filter.Page)
	q.Add(constant.Size, filter.Limit)
	q.Add("keyword", filter.Keyword)
	q.Add("status", filter.Status)
	q.Add("role_id", filter.RoleID)
	q.Add(constant.StartDate, filter.StartDate)
	q.Add(constant.EndDate, filter.EndDate)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[[]*MstMember](resp)
	if err != nil {
		return nil, nil, err
	}

	return apiResp.Data, apiResp.Meta, nil
}

func (repo *repository) CallUpdateMemberAPI(id string, reqBody map[string]interface{}) error {
	url := fmt.Sprintf(`%v/api/core/member/updateprofile/%v`, repo.cfg.Env.AifcoreHost, id)

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) CallDeleteMemberAPI(id string) error {
	url := fmt.Sprintf(`%v/api/core/member/deletemember/%v`, repo.cfg.Env.AifcoreHost, id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}
