package job

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/transaction"
	"time"
)

func NewService(repo Repository, transactionRepo transaction.Repository) Service {
	return &service{
		repo,
		transactionRepo,
	}
}

type service struct {
	repo            Repository
	transactionRepo transaction.Repository
}

type Service interface {
	CreateProCatJob(req *CreateJobRequest) (*createJobRespData, error)
	UpdateJobAPI(jobId string, req *UpdateJobRequest) error
	GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error)
	ExportJobDetailToCSV(filter *logFilter, buf *bytes.Buffer) (string, error)
	GetProCatJobDetails(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error)
	ExportJobDetailsToCSV(filter *logFilter, buf *bytes.Buffer) (string, error)
	FinalizeJob(jobIdStr string, transactionId string) error
	FinalizeFailedJob(jobIdStr string) error
}

func (svc *service) CreateProCatJob(req *CreateJobRequest) (*createJobRespData, error) {
	result, err := svc.repo.CreateJobAPI(req)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.FailedCreateJob)
	}

	return result, nil
}

func (svc *service) UpdateJobAPI(jobId string, req *UpdateJobRequest) error {
	data := map[string]interface{}{}

	if req.SuccessCount != nil {
		data["success_count"] = *req.SuccessCount
	}

	if req.Status != nil {
		data["status"] = *req.Status
	}

	if req.EndAt != nil {
		data["end_at"] = *req.EndAt
	}

	err := svc.repo.UpdateJobAPI(jobId, data)
	if err != nil {
		return apperror.MapRepoError(err, "failed to update job")
	}

	return nil
}

func (svc *service) GetProCatJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	result, err := svc.repo.GetJobsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch jobs")
	}

	return result, nil
}

func (svc *service) GetProCatJobDetail(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error) {
	result, err := svc.repo.GetJobDetailAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch job detail")
	}

	return result, nil
}

func (svc *service) GetProCatJobDetails(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error) {
	result, err := svc.repo.GetJobDetailsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch job detail")
	}

	return result, nil
}

func (svc *service) ExportJobDetailToCSV(filter *logFilter, buf *bytes.Buffer) (string, error) {
	var allDetails []*logTransProductCatalog
	page := 1
	pageSize := 100

	for {
		filter.Page = helper.ConvertUintToString(uint(page))
		filter.Size = helper.ConvertUintToString(uint(pageSize))

		resp, err := svc.repo.GetJobDetailAPI(filter)
		if err != nil {
			return "", apperror.MapRepoError(err, "failed to fetch job details")
		}

		if resp.Data == nil || len(resp.Data.JobDetails) == 0 {
			break
		}

		allDetails = append(allDetails, resp.Data.JobDetails...)

		if len(resp.Data.JobDetails) < pageSize {
			break
		}

		page++
	}

	headers := []string{}
	var mapper func(*logTransProductCatalog) []string

	switch filter.ProductSlug {
	case constant.SlugLoanRecordChecker:
		headers = []string{"Name", "NIK", "Phone Number", "Remarks", "Status", "Description"}
		mapper = mapLoanRecordCheckerRow
	}

	err := writeToCSV[*logTransProductCatalog](buf, headers, allDetails, mapper)
	if err != nil {
		return "", apperror.Internal("failed to write CSV", err)
	}

	filename := formatCSVFileName("job_detail", filter.StartDate, filter.EndDate, filter.JobId)
	return filename, nil
}

func (svc *service) ExportJobDetailsToCSV(filter *logFilter, buf *bytes.Buffer) (string, error) {
	var allDetails []*logTransProductCatalog
	page := 1
	pageSize := 100

	for {
		filter.Page = helper.ConvertUintToString(uint(page))
		filter.Size = helper.ConvertUintToString(uint(pageSize))

		resp, err := svc.repo.GetJobDetailsAPI(filter)
		if err != nil {
			return "", apperror.MapRepoError(err, "failed to fetch job details")
		}

		if resp == nil || resp.Data == nil {
			break
		}

		allDetails = append(allDetails, resp.Data.JobDetails...)

		if len(resp.Data.JobDetails) < pageSize {
			break
		}

		page++
	}

	headers := []string{}
	var mapper func(*logTransProductCatalog) []string

	switch filter.ProductSlug {
	case constant.SlugLoanRecordChecker:
		headers = []string{"Name", "NIK", "Phone Number", "Remarks", "Status", "Description"}
		mapper = mapLoanRecordCheckerRow
	}

	err := writeToCSV[*logTransProductCatalog](buf, headers, allDetails, mapper)
	if err != nil {
		return "", apperror.Internal("failed to write CSV", err)
	}

	filename := formatCSVFileName("job_detail", filter.StartDate, filter.EndDate, filter.JobId)
	return filename, nil
}

func (svc *service) FinalizeJob(jobIdStr string, transactionId string) error {
	if err := svc.transactionRepo.UpdateLogTransAPI(transactionId, map[string]interface{}{
		"success": helper.BoolPtr(true),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update transaction log")
	}

	count, err := svc.transactionRepo.ProcessedLogCountAPI(jobIdStr)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get success count")
	}

	if err := svc.repo.UpdateJobAPI(jobIdStr, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.ProcessedCount)),
		"status":        helper.StringPtr(constant.JobStatusDone),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update job status")
	}

	return nil
}

func (svc *service) FinalizeFailedJob(jobIdStr string) error {
	count, err := svc.transactionRepo.ProcessedLogCountAPI(jobIdStr)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get processed count request")
	}

	if err := svc.repo.UpdateJobAPI(jobIdStr, map[string]interface{}{
		"success_count": helper.IntPtr(int(count.ProcessedCount)),
		"status":        helper.StringPtr(constant.JobStatusFailed),
		"end_at":        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, "failed to update job status")
	}

	return nil
}

func writeToCSV[T any](buf *bytes.Buffer, headers []string, data []T, mapRow func(T) []string) error {
	writer := csv.NewWriter(buf)

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, item := range data {
		row := mapRow(item)
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

func formatCSVFileName(base, startDate, endDate, jobId string) string {
	if startDate == "" {
		return fmt.Sprintf("%s_id_%s.csv", base, jobId)
	}

	if endDate != "" && endDate != startDate {
		return fmt.Sprintf("%s_%s_until_%s.csv", base, startDate, endDate)
	}

	return fmt.Sprintf("%s_%s.csv", base, startDate)
}

func mapLoanRecordCheckerRow(d *logTransProductCatalog) []string {
	desc := ""
	if d.Message != nil {
		desc = *d.Message
	}

	remarks := ""
	status := ""
	if d.Data != nil {
		remarks = d.Data.Remarks
		status = d.Data.Status
	}

	return []string{
		d.Input.Name,
		d.Input.NIK,
		d.Input.PhoneNumber,
		remarks,
		status,
		desc,
	}
}
