package log

import (
	"encoding/json"
	"front-office/app/config"
	"io"
	"net/http"
)

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{Repo: repo, Cfg: cfg}
}

type service struct {
	Repo Repository
	Cfg  *config.Config
}

type Service interface {
	GetTransactionLogs() (*AifResponse, int, error)
	GetTransactionLogsByDate(companyId, date string) (*AifResponse, int, error)
	GetTransactionLogsByRangeDate(startDate, endDate, companyId, page string) (*AifResponse, int, error)
	GetTransactionLogsByMonth(companyId, month string) (*AifResponse, int, error)
}

func (svc *service) GetTransactionLogs() (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogs()
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByDate(companyId, date string) (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogsByDate(companyId, date)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByRangeDate(startDate, endDate, companyId, page string) (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogsByRangeDate(companyId, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByMonth(companyId, month string) (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogsByMonth(companyId, month)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func parseResponse(response *http.Response) (*AifResponse, error) {
	var baseResponse *AifResponse

	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponse); err != nil {
			return nil, err
		}
	}

	return baseResponse, nil
}
