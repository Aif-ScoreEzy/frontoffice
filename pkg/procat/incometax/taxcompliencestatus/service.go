package taxcompliancestatus

import (
	"fmt"
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
	TaxComplianceStatus(apiKey, memberId, companyId string, reqBody *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceRespData], error)
	BulkTaxComplianceStatus(apiKey string, memberId, companyId uint, file *multipart.FileHeader) error
}

func (svc *service) TaxComplianceStatus(apiKey, memberId, companyId string, reqBody *taxComplianceStatusRequest) (*model.ProCatAPIResponse[taxComplianceRespData], error) {
	product, err := svc.productRepo.GetProductAPI(constant.SlugTaxComplianceStatus)
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

	result, err := svc.repo.TaxComplianceStatusAPI(apiKey, jobIdStr, reqBody)
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

func (svc *service) BulkTaxComplianceStatus(apiKey string, memberId, companyId uint, file *multipart.FileHeader) error {
	product, err := svc.productRepo.GetProductAPI(constant.SlugTaxComplianceStatus)
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

	var taxComplianceReqs []*taxComplianceStatusRequest
	for i, record := range records {
		if i == 0 {
			continue
		}

		taxComplianceReqs = append(taxComplianceReqs, &taxComplianceStatusRequest{
			Npwp: record[0],
		})
	}

	var (
		wg         sync.WaitGroup
		errChan    = make(chan error, len(taxComplianceReqs))
		batchCount = 0
	)

	for _, req := range taxComplianceReqs {
		wg.Add(1)

		go func(taxComplianceReq *taxComplianceStatusRequest) {
			defer wg.Done()

			if err := svc.processTaxComplianceStatus(&taxComplianceContext{
				APIKey:         apiKey,
				JobIdStr:       jobIdStr,
				MemberIdStr:    memberIdStr,
				CompanyIdStr:   companyIdStr,
				MemberId:       memberId,
				CompanyId:      companyId,
				ProductId:      product.ProductId,
				ProductGroupId: product.ProductGroupId,
				JobId:          jobRes.JobId,
				Request:        taxComplianceReq,
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
		log.Error().Err(err).Msg("error during bulk tax compliance status prrocessing")
	}

	return svc.finalizeJob(jobIdStr)
}

func (svc *service) processTaxComplianceStatus(params *taxComplianceContext) error {
	fmt.Println("here")
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

	result, err := svc.repo.TaxComplianceStatusAPI(
		params.APIKey,
		params.JobIdStr,
		params.Request,
	)

	if err != nil {
		if err := svc.jobService.FinalizeFailedJob(params.JobIdStr); err != nil {

			return err
		}

		// var apiErr *apperror.ExternalAPIError
		// if errors.As(err, &apiErr) {

		// }

		return apperror.Internal("failed to process tax compliance status", err)
	}

	if err := svc.transactionRepo.UpdateLogTransAPI(result.TransactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update log transaction")
	}

	return nil
}

func (svc *service) finalizeJob(jobId string) error {
	count, err := svc.transactionRepo.ProcessedLogCountAPI(jobId)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get processed count request")
	}

	if err := svc.jobRepo.UpdateJobAPI(jobId, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.ProcessedCount)),
		"status":        helper.StringPtr(constant.JobStatusDone),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveJob)
	}

	return nil
}
