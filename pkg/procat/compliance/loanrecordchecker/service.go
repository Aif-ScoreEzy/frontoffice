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
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	logger "github.com/rs/zerolog/log"
	"github.com/usepzaka/validator"
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
	BulkLoanRecordChecker(apiKey string, memberId, companyId uint, file *multipart.FileHeader) error
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

func (svc *service) BulkLoanRecordChecker(apiKey string, memberId, companyId uint, file *multipart.FileHeader) error {
	product, err := svc.productRepo.CallGetProductBySlug(constant.SlugLoanRecordChecker)
	if err != nil {
		return apperror.MapRepoError(err, "failed to fetch product")
	}
	if product.ProductId == 0 {
		return apperror.NotFound("product not found")
	}

	if err := helper.ValidateUploadedFile(file, 30*1024*1024, []string{".csv"}); err != nil {
		return apperror.BadRequest(err.Error())
	}

	records, err := helper.ParseCSVFile(file, []string{"name", "nik", "phone_number"})
	if err != nil {
		return apperror.Internal("failed to parse csv", err)
	}

	memberIdStr := strconv.Itoa(int(memberId))
	companyIdStr := strconv.Itoa(int(companyId))
	jobRes, err := svc.logRepo.CallCreateProCatJobAPI(&log.CreateJobRequest{
		ProductId: product.ProductId,
		MemberId:  memberIdStr,
		CompanyId: companyIdStr,
		Total:     len(records) - 1,
	})
	if err != nil {
		return apperror.MapRepoError(err, "failed to create job")
	}
	jobIdStr := helper.ConvertUintToString(jobRes.JobId)

	var loanCheckerReqs []*loanRecordCheckerRequest
	for i, rec := range records {
		if i == 0 {
			continue
		} // skip header
		loanCheckerReqs = append(loanCheckerReqs, &loanRecordCheckerRequest{
			Name: rec[0], Nik: rec[1], Phone: rec[2],
		})
	}

	var (
		wg         sync.WaitGroup
		errChan    = make(chan error, len(loanCheckerReqs))
		batchCount = 0
	)

	for _, req := range loanCheckerReqs {
		wg.Add(1)

		go func(loanCheckerReq *loanRecordCheckerRequest) {
			defer wg.Done()

			if err := svc.validateAndProcessLoanChecker(apiKey, jobIdStr, memberIdStr, companyIdStr, memberId, companyId, product.ProductId, product.ProductGroupId, jobRes.JobId, loanCheckerReq); err != nil {
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
		logger.Error().Err(err).Msg("error during bulk loan record checker processing")
	}

	return svc.finalizeJob(jobIdStr)
}

func (svc *service) validateAndProcessLoanChecker(apiKey, jobIdStr, memberIdStr, companyIdStr string, memberId, companyId, productId, productGroupId, jobId uint, req *loanRecordCheckerRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		_ = svc.transactionRepo.CallCreateLogTransAPI(&transaction.LogTransProCatRequest{
			MemberID:       memberId,
			CompanyID:      companyId,
			ProductID:      productId,
			ProductGroupID: productGroupId,
			JobID:          jobId,
			Message:        err.Error(),
			Status:         http.StatusBadRequest,
			Success:        false,
			ResponseBody:   nil,
			Data:           nil,
			RequestBody:    req,
		})

		return apperror.BadRequest(err.Error())
	}

	result, err := svc.repo.CallLoanRecordCheckerAPI(apiKey, jobIdStr, memberIdStr, companyIdStr, req)
	if err != nil {
		if err := svc.logService.FinalizeFailedJob(jobIdStr); err != nil {
			return err
		}

		var apiErr *apperror.ExternalAPIError
		if errors.As(err, &apiErr) {
			return apperror.MapLoanError(apiErr)
		}

		return apperror.Internal("failed to process loan record checker", err)
	}

	if err := svc.transactionRepo.CallUpdateLogTransAPI(result.TransactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update transaction log")
	}

	return nil
}

func (svc *service) finalizeJob(jobId string) error {
	count, err := svc.transactionRepo.CallProcessedLogCount(jobId)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get processed count request")
	}

	if err := svc.logRepo.CallUpdateJobAPI(jobId, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.ProcessedCount)),
		"status":        helper.StringPtr(constant.JobStatusDone),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveJob)
	}

	return nil
}
