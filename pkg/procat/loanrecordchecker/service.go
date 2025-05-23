package loanrecordchecker

import (
	"encoding/json"
	"front-office/app/config"
	"io"
	"net/http"
)

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{
		Cfg:  cfg,
		Repo: repo,
	}
}

type service struct {
	Cfg  *config.Config
	Repo Repository
}

type Service interface {
	CallLoanRecordChecker(request *LoanRecordCheckerRequest, apiKey string) (*LoanRecordCheckerRawResponse, error)
}

func (svc *service) CallLoanRecordChecker(request *LoanRecordCheckerRequest, apiKey string) (*LoanRecordCheckerRawResponse, error) {
	response, err := svc.Repo.CallLoanRecordChecker(request, apiKey)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseResponse(response *http.Response) (*LoanRecordCheckerRawResponse, error) {
	var baseResponse *LoanRecordCheckerRawResponse

	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponse); err != nil {
			return nil, err
		}

		baseResponse.StatusCode = response.StatusCode
	}

	return baseResponse, nil
}
