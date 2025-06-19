package loanrecordchecker

import (
	"front-office/common/model"
	"front-office/helper"
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repository
}

type Service interface {
	LoanRecordChecker(request *LoanRecordCheckerRequest, apiKey, jobId, memberId, companyId string) (*model.ProCatAPIResponse[dataLoanRecord], error)
}

func (svc *service) LoanRecordChecker(request *LoanRecordCheckerRequest, apiKey, jobId, memberId, companyId string) (*model.ProCatAPIResponse[dataLoanRecord], error) {
	response, err := svc.repo.CallLoanRecordCheckerAPI(request, apiKey, jobId, memberId, companyId)
	if err != nil {
		return nil, err
	}

	result, err := helper.ParseProCatAPIResponse[dataLoanRecord](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
