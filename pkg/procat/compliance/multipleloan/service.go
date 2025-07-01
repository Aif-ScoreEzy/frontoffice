package multipleloan

import (
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/procat/job"
)

func NewService(
	repo Repository,
	productRepo product.Repository,
	jobRepo job.Repository,
	transactionRepo transaction.Repository,
	jobService job.Service,
) Service {
	return &service{
		repo,
		productRepo,
		jobRepo,
		transactionRepo,
		jobService,
	}
}

type service struct {
	repo            Repository
	productRepo     product.Repository
	jobRepo         job.Repository
	transactionRepo transaction.Repository
	jobService      job.Service
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

	jobRes, err := svc.jobRepo.CallCreateProCatJobAPI(&job.CreateJobRequest{
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

	result, err := loanHandler(apiKey, jobIdStr, memberId, companyId, reqBody)
	if err != nil {
		if err := svc.jobService.FinalizeFailedJob(jobIdStr); err != nil {
			return nil, err
		}

		return nil, apperror.MapRepoError(err, "failed to process multiple loan request")
	}

	if err := svc.jobService.FinalizeJob(jobIdStr, result.TransactionId); err != nil {
		return nil, err
	}

	return result, nil
}
