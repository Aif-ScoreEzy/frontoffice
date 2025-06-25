package loanrecordchecker

import (
	"errors"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/procat/log"
)

func NewService(
	repo Repository,
	productRepo product.Repository,
	logRepo log.Repository,
	transactionRepo transaction.Repository,
	logService log.Service,
) Service {
	return &service{
		repo,
		productRepo,
		logRepo,
		transactionRepo,
		logService,
	}
}

type service struct {
	repo            Repository
	productRepo     product.Repository
	logRepo         log.Repository
	transactionRepo transaction.Repository
	logService      log.Service
}

type Service interface {
	LoanRecordChecker(apiKey, memberId, companyId string, reqBody *loanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error)
}

func (svc *service) LoanRecordChecker(apiKey, memberId, companyId string, reqBody *loanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error) {
	product, err := svc.productRepo.CallGetProductBySlug(constant.SlugLoanRecordChecker)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch product")
	}
	if product.ProductId == 0 {
		return nil, apperror.NotFound("product not found")
	}

	jobRes, err := svc.logRepo.CallCreateProCatJobAPI(&log.CreateJobRequest{
		ProductId: product.ProductId,
		MemberId:  memberId,
		CompanyId: companyId,
		Total:     1,
	})
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to create job")
	}
	jobIdStr := helper.ConvertUintToString(jobRes.JobId)

	result, err := svc.repo.CallLoanRecordCheckerAPI(apiKey, jobIdStr, memberId, companyId, reqBody)
	if err != nil {
		if err := svc.logService.FinalizeFailedJob(jobIdStr); err != nil {
			return nil, err
		}

		var apiErr *apperror.ExternalAPIError
		if errors.As(err, &apiErr) {
			return nil, apperror.MapLoanError(apiErr)
		}

		return nil, apperror.Internal("failed to process loan record checker", err)
	}

	if err := svc.logService.FinalizeJob(jobIdStr, result.TransactionId); err != nil {
		return nil, err
	}

	return result, nil
}
