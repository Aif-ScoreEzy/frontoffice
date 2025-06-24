package taxscore

import (
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
	CallTaxScore(apiKey, memberId, companyId string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreDataResponse], error)
}

func (svc *service) CallTaxScore(apiKey, memberId, companyId string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreDataResponse], error) {
	product, err := svc.productRepo.CallGetProductBySlug(constant.SlugTaxComplianceStatus)
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

	result, err := svc.repo.CallTaxScoreAPI(apiKey, jobIdStr, request)
	if err != nil {
		if err := svc.logService.FinalizeFailedJob(jobIdStr); err != nil {
			return nil, err
		}

		return nil, apperror.MapRepoError(err, "failed to process tax score")
	}

	if err := svc.logService.FinalizeJob(jobIdStr, result.TransactionId); err != nil {
		return nil, err
	}

	return result, nil
}
