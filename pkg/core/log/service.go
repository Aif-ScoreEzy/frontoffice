package log

import (
	"encoding/json"
	"front-office/app/config"
	"io"
)

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{Repo: repo, Cfg: cfg}
}

type service struct {
	Repo Repository
	Cfg *config.Config
}

type Service interface {
	GetTransactionLogsSvc() (*AifResponse, int, error)
	GetTransactionLogsByDateSvc(companyId, date string) (*AifResponse, int, error)
	GetTransactionLogsByRangeDateSvc(startDate, endDate, companyId, page string) (*AifResponse, int, error)
	GetTransactionLogsByMonthSvc(companyId, month string) (*AifResponse, int, error)
}

func (svc *service) GetTransactionLogsSvc() (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogs()
	if err != nil {
		return nil, 0, err
	}

	var baseResponse *AifResponse
	if response != nil {
		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, 0, err
		}
		
		defer response.Body.Close()
		
		if err := json.Unmarshal(responseBodyBytes, &baseResponse); err != nil {
			return nil, 0, err
		}
	}

	return baseResponse, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByDateSvc(companyId, date string) (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogsByDate(companyId, date)
	if err != nil {
		return nil, 0, err
	}

	var baseResponse *AifResponse
	if response != nil {
		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, 0, err
		}
		
		defer response.Body.Close()
		
		if err := json.Unmarshal(responseBodyBytes, &baseResponse); err != nil {
			return nil, 0, err
		}
	}

	return baseResponse, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByRangeDateSvc(startDate, endDate, companyId, page string) (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogsByRangeDate(companyId, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}

	var baseResponse *AifResponse
	if response != nil {
		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, 0, err
		}
		
		defer response.Body.Close()
		
		if err := json.Unmarshal(responseBodyBytes, &baseResponse); err != nil {
			return nil, 0, err
		}
	}

	return baseResponse, response.StatusCode, nil
}

func (svc *service) GetTransactionLogsByMonthSvc(companyId, month string) (*AifResponse, int, error) {
	response, err := svc.Repo.FindAllTransactionLogsByMonth(companyId, month)
	if err != nil {
		return nil, 0, err
	}

	var baseResponse *AifResponse
	if response != nil {
		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, 0, err
		}
		
		defer response.Body.Close()
		
		if err := json.Unmarshal(responseBodyBytes, &baseResponse); err != nil {
			return nil, 0, err
		}
	}

	return baseResponse, response.StatusCode, nil
}
