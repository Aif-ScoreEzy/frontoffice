package taxscore

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
	CallTaxScore(apiKey, memberId, companyId string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreRespData], error)
}

func (svc *service) CallTaxScore(apiKey, memberId, companyId string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreRespData], error) {
	product, err := svc.productRepo.GetProductAPI(constant.SlugTaxScore)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedFetchProduct)
	}
	if product.ProductId == 0 {
		return nil, apperror.NotFound(constant.ProductNotFound)
	}

	jobRes, err := svc.jobRepo.CreateJobAPI(&job.CreateJobRequest{
		ProductId: product.ProductId,
		MemberId:  memberId,
		CompanyId: companyId,
		Total:     1,
	})
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedCreateJob)
	}
	jobIdStr := helper.ConvertUintToString(jobRes.JobId)

	result, err := svc.repo.CallTaxScoreAPI(apiKey, jobIdStr, request)
	if err != nil {
		if err := svc.jobService.FinalizeFailedJob(jobIdStr); err != nil {
			return nil, err
		}

		return nil, apperror.MapRepoError(err, "failed to process tax score")
	}

	if err := svc.jobService.FinalizeJob(jobIdStr, result.TransactionId); err != nil {
		return nil, err
	}

	return result, nil
}
