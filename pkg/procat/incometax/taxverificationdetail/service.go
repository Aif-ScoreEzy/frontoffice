package taxverificationdetail

import (
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/procat/job"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/usepzaka/validator"
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
	CallTaxVerification(apiKey, memberId, companyId string, request *taxVerificationRequest) (*model.ProCatAPIResponse[taxVerificationRespData], error)
	BulkTaxVerification(apiKey string, memberId, companyId uint, file *multipart.FileHeader) error
}

func (svc *service) CallTaxVerification(apiKey, memberId, companyId string, request *taxVerificationRequest) (*model.ProCatAPIResponse[taxVerificationRespData], error) {
	product, err := svc.productRepo.GetProductAPI(constant.SlugTaxVerificationDetail)
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

	result, err := svc.repo.TaxVerificationAPI(apiKey, jobIdStr, request)
	if err != nil {
		if err := svc.jobService.FinalizeFailedJob(jobIdStr); err != nil {
			return nil, err
		}

		return nil, apperror.MapRepoError(err, "failed to process tax score")
	}

	if err := svc.transactionRepo.UpdateLogTransAPI(result.TransactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return nil, apperror.MapRepoError(err, "failed to update transaction log")
	}

	if err := svc.jobService.FinalizeJob(jobIdStr); err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) BulkTaxVerification(apiKey string, memberId, companyId uint, file *multipart.FileHeader) error {
	product, err := svc.productRepo.GetProductAPI(constant.SlugTaxVerificationDetail)
	if err != nil {
		return apperror.MapRepoError(err, constant.FailedFetchProduct)
	}
	if product.ProductId == 0 {
		return apperror.NotFound(constant.ProductNotFound)
	}

	if err := helper.ValidateUploadedFile(file, 30*1024*1024, []string{".csv"}); err != nil {
		return apperror.BadRequest(err.Error())
	}

	records, err := helper.ParseCSVFile(file, []string{"npwp"})
	if err != nil {
		return apperror.Internal(constant.FailedParseCSV, err)
	}

	memberIdStr := strconv.Itoa(int(memberId))
	companyIdStr := strconv.Itoa(int(companyId))
	jobRes, err := svc.jobRepo.CreateJobAPI(&job.CreateJobRequest{
		ProductId: product.ProductId,
		MemberId:  memberIdStr,
		CompanyId: companyIdStr,
		Total:     len(records) - 1,
	})
	if err != nil {
		return apperror.MapRepoError(err, constant.FailedCreateJob)
	}
	jobIdStr := helper.ConvertUintToString(jobRes.JobId)

	var taxScoreReqs []*taxVerificationRequest
	for i, record := range records {
		if i == 0 {
			continue
		}

		taxScoreReqs = append(taxScoreReqs, &taxVerificationRequest{
			NpwpOrNik: record[0],
		})
	}

	var (
		wg         sync.WaitGroup
		errChan    = make(chan error, len(taxScoreReqs))
		batchCount = 0
	)

	for _, req := range taxScoreReqs {
		wg.Add(1)

		go func(taxScoreReq *taxVerificationRequest) {
			defer wg.Done()

			if err := svc.processTaxVerification(&taxVerificationContext{
				APIKey:         apiKey,
				JobIdStr:       jobIdStr,
				MemberIdStr:    memberIdStr,
				CompanyIdStr:   companyIdStr,
				MemberId:       memberId,
				CompanyId:      companyId,
				ProductId:      product.ProductId,
				ProductGroupId: product.ProductGroupId,
				JobId:          jobRes.JobId,
				Request:        taxScoreReq,
			}); err != nil {
				errChan <- err
			}
		}(req)

		batchCount++
		if batchCount == 100 {
			time.Sleep(time.Second)
			batchCount = 0
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Error().Err(err).Msg("error during bulk tax score prrocessing")
	}

	return svc.jobService.FinalizeJob(jobIdStr)
}

func (svc *service) processTaxVerification(params *taxVerificationContext) error {
	if err := validator.ValidateStruct(params.Request); err != nil {
		_ = svc.transactionRepo.CreateLogTransAPI(&transaction.LogTransProCatRequest{
			MemberID:       params.MemberId,
			CompanyID:      params.CompanyId,
			ProductID:      params.ProductId,
			ProductGroupID: params.ProductGroupId,
			JobID:          params.JobId,
			Success:        false,
			Message:        err.Error(),
			Status:         http.StatusBadRequest,
			ResponseBody:   nil,
			Data:           nil,
			RequestBody:    params.Request,
		})

		return apperror.BadRequest(err.Error())
	}

	result, err := svc.repo.TaxVerificationAPI(
		params.APIKey,
		params.JobIdStr,
		params.Request,
	)

	if err != nil {
		if err := svc.jobService.FinalizeFailedJob(params.JobIdStr); err != nil {

			return err
		}

		return apperror.Internal("failed to process tax compliance status", err)
	}

	if err := svc.transactionRepo.UpdateLogTransAPI(result.TransactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update log transaction")
	}

	return nil
}
