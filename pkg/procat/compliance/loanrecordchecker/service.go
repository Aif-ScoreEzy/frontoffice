package loanrecordchecker

import (
	"front-office/app/config"
	"front-office/common/model"
	"front-office/helper"
	"log"
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
	LoanRecordChecker(request *LoanRecordCheckerRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataLoanRecord], error)
}

func (svc *service) LoanRecordChecker(request *LoanRecordCheckerRequest, apiKey, memberId, companyId string) (*model.ProCatAPIResponse[dataLoanRecord], error) {
	response, err := svc.Repo.CallLoanRecordCheckerAPI(request, apiKey, memberId, companyId)
	log.Println("loan record responseee====>", response)

	if err != nil {
		return nil, err
	}

	result, err := helper.ParseProCatAPIResponse[dataLoanRecord](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
