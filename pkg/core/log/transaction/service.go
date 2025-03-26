package transaction

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
	GetLogTransactions() (*AifResponse, int, error)
	GetLogTransactionsByDate(companyId, date string) (*AifResponse, int, error)
	GetLogTransactionsByRangeDate(startDate, endDate, companyId, page string) (*AifResponse, int, error)
	GetLogTransactionsByMonth(companyId, month string) (*AifResponse, int, error)
}

func (svc *service) GetLogTransactions() (*AifResponse, int, error) {
	response, err := svc.Repo.FetchLogTransactions()
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetLogTransactionsByDate(companyId, date string) (*AifResponse, int, error) {
	response, err := svc.Repo.FetchLogTransactionsByDate(companyId, date)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetLogTransactionsByRangeDate(startDate, endDate, companyId, page string) (*AifResponse, int, error) {
	response, err := svc.Repo.FetchLogTransactionsByRangeDate(companyId, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetLogTransactionsByMonth(companyId, month string) (*AifResponse, int, error) {
	response, err := svc.Repo.FetchLogTransactionsByMonth(companyId, month)
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
