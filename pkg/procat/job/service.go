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
	"strconv"
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
	CreateJob(req *CreateJobRequest) (*createJobRespData, error)
	UpdateJobAPI(jobId string, req *UpdateJobRequest) error
	GetJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error)
	GetJobDetails(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error)
	ExportJobDetails(filter *logFilter, buf *bytes.Buffer) (string, error)
	GetJobDetailsByDateRange(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error)
	ExportJobDetailsByDateRange(filter *logFilter, buf *bytes.Buffer) (string, error)
	FinalizeJob(jobIdStr string) error
	FinalizeFailedJob(jobIdStr string) error
}

func (svc *service) CreateJob(req *CreateJobRequest) (*createJobRespData, error) {
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

func (svc *service) GetJob(filter *logFilter) (*model.AifcoreAPIResponse[any], error) {
	result, err := svc.repo.GetJobsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch jobs")
	}

	return result, nil
}

func (svc *service) GetJobDetails(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error) {
	result, err := svc.repo.GetJobDetailAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch job detail")
	}

	return result, nil
}

func (svc *service) GetJobDetailsByDateRange(filter *logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error) {
	result, err := svc.repo.GetJobDetailsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch job detail")
	}

	return result, nil
}

func (svc *service) ExportJobDetails(filter *logFilter, buf *bytes.Buffer) (string, error) {
	return svc.exportJobDetailsToCSV(filter, buf, svc.repo.GetJobDetailAPI, false)
}

func (svc *service) ExportJobDetailsByDateRange(filter *logFilter, buf *bytes.Buffer) (string, error) {
	return svc.exportJobDetailsToCSV(filter, buf, svc.repo.GetJobDetailsAPI, true)
}

func (svc *service) exportJobDetailsToCSV(
	filter *logFilter,
	buf *bytes.Buffer,
	fetchFunc func(*logFilter) (*model.AifcoreAPIResponse[*jobDetailResponse], error),
	includeDate bool,
) (string, error) {
	resp, err := fetchFunc(filter)
	if err != nil {
		return "", apperror.MapRepoError(err, "failed to fetch job details")
	}

	headers := []string{}
	var mapper func(*logTransProductCatalog) []string

	switch filter.ProductSlug {
	case constant.SlugLoanRecordChecker:
		headers = []string{"Name", "NIK", "Phone Number", "Remarks", "Data Status", "Status", "Description"}
		mapper = mapLoanRecordCheckerRow
	case constant.SlugMultipleLoan7Days, constant.SlugMultipleLoan30Days, constant.SlugMultipleLoan90Days:
		headers = []string{"NIK", "Phone Number", "Query Count", "Status", "Description"}
		mapper = mapMultipleLoanRow
	case constant.SlugTaxComplianceStatus:
		headers = []string{"NPWP", "Nama", "Alamat", "Data Status", "Status", "Description"}
		mapper = mapTaxComplianceRow
	case constant.SlugTaxScore:
		headers = []string{"NPWP", "Nama", "Alamat", "Data Status", "Score", "Status", "Description"}
		mapper = mapTaxScoreRow
	case constant.SlugTaxVerificationDetail:
		headers = []string{"NPWP Or NIK", "Nama", "Alamat", "NPWP", "NPWP Verification", "Data Status", "Tax Compliance", "Status", "Description"}
		mapper = mapTaxVerificationRow
	}

	if includeDate {
		headers = append([]string{"Date"}, headers...)
		mapper = withDateColumn(mapper)
	}

	err = writeToCSV[*logTransProductCatalog](buf, headers, resp.Data.JobDetails, mapper)
	if err != nil {
		return "", apperror.Internal("failed to write CSV", err)
	}

	filename := formatCSVFileName("job_detail", filter.StartDate, filter.EndDate, filter.JobId)
	return filename, nil
}

func (svc *service) FinalizeJob(jobIdStr string) error {
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

type rowMapper func(*logTransProductCatalog) []string

func withDateColumn(mapper rowMapper) rowMapper {
	return func(d *logTransProductCatalog) []string {
		row := mapper(d)
		date := d.DateTime

		return append([]string{date}, row...)
	}
}

func mapLoanRecordCheckerRow(d *logTransProductCatalog) []string {
	desc := ""
	if d.Message != nil {
		desc = *d.Message
	}

	remarks := ""
	status := ""
	if d.Data != nil {
		remarks = *d.Data.Remarks
		status = *d.Data.Status
	}

	return []string{
		*d.Input.Name,
		*d.Input.NIK,
		*d.Input.PhoneNumber,
		remarks,
		status,
		d.Status,
		desc,
	}
}

func mapMultipleLoanRow(d *logTransProductCatalog) []string {
	desc := ""
	if d.Message != nil {
		desc = *d.Message
	}

	queryCount := 0
	if d.Data != nil {
		queryCount = *d.Data.QueryCount
	}

	return []string{
		*d.Input.NIK,
		*d.Input.PhoneNumber,
		strconv.Itoa(queryCount),
		d.Status,
		desc,
	}
}

func mapTaxComplianceRow(d *logTransProductCatalog) []string {
	desc := ""
	if d.Message != nil {
		desc = *d.Message
	}

	nama := ""
	alamat := ""
	status := ""

	if d.Data != nil {
		nama = *d.Data.Nama
		alamat = *d.Data.Alamat
		status = *d.Data.Status
	}

	return []string{
		*d.Input.NPWP,
		nama,
		alamat,
		status,
		d.Status,
		desc,
	}
}

func mapTaxScoreRow(d *logTransProductCatalog) []string {
	desc := ""
	if d.Message != nil {
		desc = *d.Message
	}

	nama := ""
	alamat := ""
	status := ""
	score := ""

	if d.Data != nil {
		nama = *d.Data.Nama
		alamat = *d.Data.Alamat
		status = *d.Data.Status
		score = *d.Data.Score
	}

	return []string{
		*d.Input.NPWP,
		nama,
		alamat,
		status,
		score,
		d.Status,
		desc,
	}
}

func mapTaxVerificationRow(d *logTransProductCatalog) []string {
	desc := ""
	if d.Message != nil {
		desc = *d.Message
	}

	nama := ""
	alamat := ""
	npwp := ""
	npwpVerification := ""
	status := ""
	taxCompliance := ""

	if d.Data != nil {
		nama = *d.Data.Nama
		alamat = *d.Data.Alamat
		npwp = *d.Data.NPWP
		npwpVerification = *d.Data.NPWPVerification
		status = *d.Data.Status
		taxCompliance = *d.Data.TaxCompliance
	}

	return []string{
		*d.Input.NPWPOrNIK,
		nama,
		alamat,
		npwp,
		npwpVerification,
		status,
		taxCompliance,
		d.Status,
		desc,
	}
}
