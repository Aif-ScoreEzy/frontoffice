package taxcompliancestatus

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
	CallTaxCompliance(apiKey, memberId, companyId string, reqBody *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceRespData], error)
}

func (svc *service) CallTaxCompliance(apiKey, memberId, companyId string, reqBody *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceRespData], error) {
	product, err := svc.productRepo.GetProductAPI(constant.SlugTaxComplianceStatus)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch product")
	}
	if product.ProductId == 0 {
		return nil, apperror.NotFound("product not found")
	}

	jobRes, err := svc.jobRepo.CallCreateJobAPI(&job.CreateJobRequest{
		ProductId: product.ProductId,
		MemberId:  memberId,
		CompanyId: companyId,
		Total:     1,
	})
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to create job")
	}
	jobIdStr := helper.ConvertUintToString(jobRes.JobId)

	result, err := svc.repo.CallTaxComplianceStatusAPI(apiKey, jobIdStr, reqBody)
	if err != nil {
		if err := svc.jobService.FinalizeFailedJob(jobIdStr); err != nil {
			return nil, err
		}

		return nil, apperror.MapRepoError(err, "failed to process tax compliance status")
	}

	if err := svc.jobService.FinalizeJob(jobIdStr, result.TransactionId); err != nil {
		return nil, err
	}

	return result, nil
}
