package log

import (
	"encoding/json"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"io"
	"net/http"
)

func NewService(cfg *config.Config) Service {
	return &service{Cfg: cfg}
}

type service struct {
	Cfg *config.Config
}

type Service interface {
	GetTransactionLogsByDateSvc(companyId, date string) (*model.AifResponse, int, error)
	GetTransactionLogsByRangeDateSvc(startDate, endDate, companyId, page string) (*model.AifResponse, int, error)
	GetTransactionLogsByMonthSvc(companyId, month string) (*model.AifResponse, int, error)
	GetTransactionLogsByNameSvc(companyId, name string) (*model.AifResponse, int, error)
}

func (svc *service) GetTransactionLogsByDateSvc(companyId, date string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/by"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("date", date)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResp)

	return dataResp, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByRangeDateSvc(startDate, endDate, companyId, page string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/range"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("date_start", startDate)
	q.Add("date_end", endDate)
	q.Add("company_id", companyId)
	q.Add("page", page)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResp)

	return dataResp, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByMonthSvc(companyId, month string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/core/logging/transaction/month"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("month", month)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResp)

	return dataResp, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByNameSvc(companyId, name string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/log/byname"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyId)
	q.Add("name", name)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResp)

	return dataResp, response.StatusCode, nil
}
