package loanrecordchecker

import (
	"encoding/json"
	"front-office/app/config"
	"front-office/common/model"
	"front-office/helper"
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
	LoanRecordChecker(request *LoanRecordCheckerRequest, apiKey, memberId, companyId string) (*LoanRecordCheckerRawResponse, error)
	GetLoanRecordCheckerJob(filter *loanRecordCheckerFilter) (*model.AifcoreAPIResponse[any], error)
	GetLoanRecordCheckerJobDetail(filter *loanRecordCheckerFilter) (*model.AifcoreAPIResponse[any], error)
}

func (svc *service) LoanRecordChecker(request *LoanRecordCheckerRequest, apiKey, memberId, companyId string) (*LoanRecordCheckerRawResponse, error) {
	response, err := svc.Repo.CallLoanRecordCheckerAPI(request, apiKey, memberId, companyId)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) GetLoanRecordCheckerJob(filter *loanRecordCheckerFilter) (*model.AifcoreAPIResponse[any], error) {
	response, err := svc.Repo.CallGetLoanRecordCheckerJobAPI(filter)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}

func (svc *service) GetLoanRecordCheckerJobDetail(filter *loanRecordCheckerFilter) (*model.AifcoreAPIResponse[any], error) {
	response, err := svc.Repo.CallGetLoanRecordCheckerJobDetailAPI(filter)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
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
