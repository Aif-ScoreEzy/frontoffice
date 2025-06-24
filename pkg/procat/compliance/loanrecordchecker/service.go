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
	"time"
)

func NewService(
	repo Repository,
	productRepo product.Repository,
	logRepo log.Repository,
	transactionRepo transaction.Repository,
) Service {
	return &service{
		repo,
		productRepo,
		logRepo,
		transactionRepo,
	}
}

type service struct {
	repo            Repository
	productRepo     product.Repository
	logRepo         log.Repository
	transactionRepo transaction.Repository
}

type Service interface {
	LoanRecordChecker(apiKey, memberId, companyId string, reqBody *LoanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error)
}

func (svc *service) LoanRecordChecker(apiKey, memberId, companyId string, reqBody *LoanRecordCheckerRequest) (*model.ProCatAPIResponse[dataLoanRecord], error) {
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
		var apiErr *apperror.ExternalAPIError
		if errors.As(err, &apiErr) {
			return nil, apperror.MapLoanError(apiErr)
		}

		return nil, apperror.Internal("failed to process loan record checker", err)
	}

	// mark success on transaction log
	if err := svc.transactionRepo.CallUpdateLogTransAPI(result.TransactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return nil, apperror.MapRepoError(err, "failed to update log")
	}

	// get count of success
	count, err := svc.transactionRepo.CallLogTransSuccessCountAPI(jobIdStr)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to get success count")
	}

	if err := svc.logRepo.CallUpdateJobAPI(jobIdStr, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.SuccessCount)),
		"status":        helper.StringPtr(constant.JobStatusDone),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return nil, apperror.MapRepoError(err, "failed to update job")
	}

	return result, nil
}
