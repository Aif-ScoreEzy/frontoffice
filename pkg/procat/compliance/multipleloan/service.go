package multipleloan

import (
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
	MultipleLoan(apiKey, productSlug, memberId, companyId string, reqBody *multipleLoanRequest) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error)
}

type multipleLoanFunc func(string, string, string, string, *multipleLoanRequest) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error)

func (svc *service) MultipleLoan(apiKey, productSlug, memberId, companyId string, reqBody *multipleLoanRequest) (*model.ProCatAPIResponse[dataMultipleLoanResponse], error) {
	product, err := svc.productRepo.CallGetProductBySlug(productSlug)
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

	loanHandlers := map[string]multipleLoanFunc{
		constant.SlugMultipleLoan7Days:  svc.repo.CallMultipleLoan7Days,
		constant.SlugMultipleLoan30Days: svc.repo.CallMultipleLoan30Days,
		constant.SlugMultipleLoan90Days: svc.repo.CallMultipleLoan90Days,
	}

	loanHandler, ok := loanHandlers[productSlug]
	if !ok {
		return nil, apperror.BadRequest("unsupported product type")
	}

	result, err := loanHandler(apiKey, memberId, jobIdStr, companyId, reqBody)
	if err != nil {
		if err := svc.logRepo.CallUpdateJobAPI(jobIdStr, map[string]interface{}{
			"success_count": helper.IntPtr(0),
			"status":        helper.StringPtr(constant.JobStatusFailed),
			"end_at":        helper.TimePtr(time.Now()),
		}); err != nil {
			return nil, apperror.MapRepoError(err, "failed to update job status")
		}

		return nil, apperror.MapRepoError(err, "failed to process multiple loan request")
	}

	if err := svc.logService.FinalizeLoanJob(jobIdStr, result.TransactionId); err != nil {
		return nil, err
	}

	return result, nil
}
