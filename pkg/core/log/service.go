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
	GetTransactionLogsByDateSvc(companyID, date string) (*model.AifResponse, int, error)
	GetTransactionLogsByRangeDateSvc(startDate, endDate, companyID, page string) (*model.AifResponse, int, error)
	GetTransactionLogsByMonthSvc(companyID, month string) (*model.AifResponse, int, error)
	GetTransactionLogsByNameSvc(companyID, name string) (*model.AifResponse, int, error)
}

func (svc *service) GetTransactionLogsByDateSvc(companyID, date string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/log/by"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyID)
	q.Add("date", date)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err := json.Unmarshal(responseBodyBytes, &dataResp); err != nil {
		return nil, 0, err
	}

	return dataResp, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByRangeDateSvc(startDate, endDate, companyID, page string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/log/byrange"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("date_start", startDate)
	q.Add("date_end", endDate)
	q.Add("company_id", companyID)
	q.Add("page", page)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err := json.Unmarshal(responseBodyBytes, &dataResp); err != nil {
		return nil, 0, err
	}

	return dataResp, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByMonthSvc(companyID, month string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/log/bymonth"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyID)
	q.Add("month", month)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err := json.Unmarshal(responseBodyBytes, &dataResp); err != nil {
		return nil, 0, err
	}

	return dataResp, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByNameSvc(companyID, name string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := svc.Cfg.Env.AifcoreHost + "/api/log/byname"

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyID)
	q.Add("name", name)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err := json.Unmarshal(responseBodyBytes, &dataResp); err != nil {
		return nil, 0, err
	}

	return dataResp, response.StatusCode, nil
}
