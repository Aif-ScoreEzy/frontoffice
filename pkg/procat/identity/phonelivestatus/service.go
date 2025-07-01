package phonelivestatus

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/member"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/usepzaka/validator"
)

func NewService(repo Repository, memberRepo member.Repository) Service {
	return &service{
		repo,
		memberRepo,
	}
}

type service struct {
	repo       Repository
	memberRepo member.Repository
}

type Service interface {
	CreateJob(memberId, companyId string, reqBody *createJobRequest) (*createJobRespData, error)
	GetPhoneLiveStatusJob(filter *phoneLiveStatusFilter) (*jobListRespData, error)
	GetAllPhoneLiveStatusDetails(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	GetPhoneLiveStatusDetailsByRangeDate(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	GetJobsSummary(filter *phoneLiveStatusFilter) (*jobsSummaryRespData, error)
	GetPhoneLiveStatusDetailsSummary(filter *phoneLiveStatusFilter) (*jobDetailRespData, error)
	ExportJobsSummary(filter *phoneLiveStatusFilter, buf *bytes.Buffer) (string, error)
	ExportJobDetails(filter *phoneLiveStatusFilter, buf *bytes.Buffer) (string, error)
	UpdateJob(jobId uint, reqBody *updateJobRequest) error
	UpdateJobDetail(jobId, jobDetailId uint, reqBody *updateJobDetailRequest) error
	ProcessPhoneLiveStatus(memberId, companyId string, reqBody *phoneLiveStatusRequest) error
	BulkProcessPhoneLiveStatus(apiKey, memberId, companyId string, fileHeader *multipart.FileHeader) error
}

func (svc *service) CreateJob(memberId, companyId string, reqBody *createJobRequest) (*createJobRespData, error) {
	job, err := svc.repo.CallCreateJobAPI(reqBody)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.ErrCreatePhoneLiveJob)
	}

	return job, err
}

func (svc *service) GetPhoneLiveStatusJob(filter *phoneLiveStatusFilter) (*jobListRespData, error) {
	jobs, err := svc.repo.CallGetPhoneLiveStatusJobAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch phone live status jobs")
	}

	return jobs, nil
}

func (svc *service) GetPhoneLiveStatusDetailsSummary(filter *phoneLiveStatusFilter) (*jobDetailRespData, error) {
	jobDetails, err := svc.repo.CallGetJobDetailsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.ErrFetchPhoneLiveDetail)
	}

	return jobDetails, nil
}

func (svc *service) GetAllPhoneLiveStatusDetails(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error) {
	jobDetails, err := svc.repo.CallGetAllJobDetailsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.ErrFetchPhoneLiveDetail)
	}

	return jobDetails, nil
}

func (svc *service) GetPhoneLiveStatusDetailsByRangeDate(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error) {
	jobDetail, err := svc.repo.CallGetJobDetailsByRangeDateAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, constant.ErrFetchPhoneLiveDetail)
	}

	return jobDetail, nil
}

func (svc *service) GetJobsSummary(filter *phoneLiveStatusFilter) (*jobsSummaryRespData, error) {
	jobsSummary, err := svc.repo.CallGetJobsSummary(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch phone live status jobs summary")
	}

	return jobsSummary, nil
}

func (svc *service) ExportJobsSummary(filter *phoneLiveStatusFilter, buf *bytes.Buffer) (string, error) {
	data, err := svc.repo.CallGetJobDetailsByRangeDateAPI(filter)
	if err != nil {
		return "", apperror.MapRepoError(err, "failed to fetch job details")
	}

	if err := writeJobDetailsToCSV(buf, data); err != nil {
		return "", apperror.Internal("failed to write CSV", err)
	}

	filename := formatCSVFileName("job_summary", filter.StartDate, filter.EndDate)

	return filename, nil
}

func (svc *service) ExportJobDetails(filter *phoneLiveStatusFilter, buf *bytes.Buffer) (string, error) {
	data, err := svc.repo.CallGetAllJobDetailsAPI(filter)
	if err != nil {
		return "", apperror.MapRepoError(err, "failed to fetch job details")
	}

	if err := writeJobDetailsToCSV(buf, data); err != nil {
		return "", apperror.Internal("failed to write CSV", err)
	}

	filename := formatCSVFileName("job_summary", filter.StartDate, filter.EndDate)

	return filename, nil
}

func (svc *service) UpdateJob(jobId uint, reqBody *updateJobRequest) error {
	jobIdStr := strconv.FormatUint(uint64(jobId), 10)
	if err := svc.repo.CallUpdateJob(jobIdStr, reqBody); err != nil {
		return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveJob)
	}

	return nil
}

func (svc *service) UpdateJobDetail(jobId, jobDetailId uint, reqBody *updateJobDetailRequest) error {
	jobIdStr := strconv.FormatUint(uint64(jobId), 10)
	jobDetailIdStr := strconv.FormatUint(uint64(jobId), 10)
	if err := svc.repo.CallUpdateJobDetail(jobIdStr, jobDetailIdStr, reqBody); err != nil {
		return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveDetail)
	}

	return nil
}

func (svc *service) ProcessPhoneLiveStatus(memberId, companyId string, reqBody *phoneLiveStatusRequest) error {
	member, err := svc.memberRepo.CallGetMemberAPI(&member.FindUserQuery{
		Id:        memberId,
		CompanyId: companyId,
	})
	if err != nil {
		return apperror.MapRepoError(err, constant.FailedFetchMember)
	}
	if member.MemberId == 0 {
		return apperror.NotFound(constant.UserNotFound)
	}

	job, err := svc.repo.CallCreateJobAPI(&createJobRequest{
		MemberId:                memberId,
		CompanyId:               companyId,
		PhoneLiveStatusRequests: []*phoneLiveStatusRequest{reqBody},
	})
	if err != nil {
		return apperror.MapRepoError(err, constant.ErrCreatePhoneLiveJob)
	}
	jobIdStr := strconv.Itoa(int(job.JobId))

	jobDetails, err := svc.repo.CallGetAllJobDetailsAPI(&phoneLiveStatusFilter{
		JobId:     jobIdStr,
		MemberId:  memberId,
		CompanyId: companyId,
	})
	if err != nil {
		return apperror.MapRepoError(err, "failed to fetch job detail")
	}
	if len(jobDetails) == 0 {
		return apperror.NotFound("no job details found")
	}
	jobDetailIdStr := strconv.Itoa(int(jobDetails[0].Id))

	if err := svc.validateSingleRequest(jobIdStr, jobDetailIdStr, reqBody); err != nil {
		return err
	}

	result, err := svc.processPhoneLiveStatus(member.Key, jobIdStr, jobDetailIdStr, jobDetails[0])
	if err != nil {
		return err
	}

	if err := svc.updateProcessedDetail(jobIdStr, jobDetailIdStr, result); err != nil {
		return err
	}

	return svc.finalizeJob(jobIdStr)
}

func (svc *service) BulkProcessPhoneLiveStatus(apiKey, memberId, companyId string, file *multipart.FileHeader) error {
	if err := helper.ValidateUploadedFile(file, 30*1024*1024, []string{".csv"}); err != nil {
		return apperror.BadRequest(err.Error())
	}

	records, err := helper.ParseCSVFile(file, []string{"phone_number"})
	if err != nil {
		return apperror.Internal("failed to parse csv", err)
	}

	var phoneReqs []*phoneLiveStatusRequest
	for i := 1; i < len(records); i++ { // Skip header
		phoneReqs = append(phoneReqs, &phoneLiveStatusRequest{
			PhoneNumber: records[i][0],
		})
	}

	job, err := svc.repo.CallCreateJobAPI(&createJobRequest{
		MemberId:                memberId,
		CompanyId:               companyId,
		PhoneLiveStatusRequests: phoneReqs,
	})
	if err != nil {
		return apperror.MapRepoError(err, constant.ErrCreatePhoneLiveJob)
	}
	jobIdStr := strconv.Itoa(int(job.JobId))

	jobDetails, err := svc.repo.CallGetAllJobDetailsAPI(&phoneLiveStatusFilter{
		JobId:     jobIdStr,
		MemberId:  memberId,
		CompanyId: companyId,
	})
	if err != nil {
		return apperror.MapRepoError(err, "failed to fetch job detail")
	}
	if len(jobDetails) == 0 {
		return apperror.NotFound("no job details found")
	}

	var (
		wg         sync.WaitGroup
		errChan    = make(chan error, len(jobDetails))
		batchCount = 0
	)

	for _, detail := range jobDetails {
		wg.Add(1)

		go func(detail *mstPhoneLiveStatusJobDetail) {
			defer wg.Done()
			detailID := strconv.Itoa(int(detail.Id))

			if err := svc.processAndUpdatePhoneLiveStatus(apiKey, jobIdStr, detailID, detail); err != nil {
				errChan <- err
			}
		}(detail)

		batchCount++
		if batchCount == 100 {
			time.Sleep(time.Second)
			batchCount = 0
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Error().Err(err).Msg("error during bulk phone live status processing")
	}

	return svc.finalizeJob(jobIdStr)
}

func writeJobDetailsToCSV(buf *bytes.Buffer, data []*mstPhoneLiveStatusJobDetail) error {
	w := csv.NewWriter(buf)
	headers := []string{"Phone Number", "Subscriber Status", "Device Status", "Status", "Operator", "Phone Type"}

	if err := w.Write(headers); err != nil {
		return err
	}

	for _, d := range data {
		row := []string{d.PhoneNumber, d.SubscriberStatus, d.DeviceStatus, d.Status, d.Operator, d.PhoneType}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func formatCSVFileName(base, startDate, endDate string) string {
	if endDate != "" && endDate != startDate {
		return fmt.Sprintf("%s_%s_until_%s.csv", base, startDate, endDate)
	}
	return fmt.Sprintf("%s_%s.csv", base, startDate)
}

func (svc *service) validateSingleRequest(jobId, jobDetailId string, reqBody *phoneLiveStatusRequest) error {
	if errValidation := validator.ValidateStruct(reqBody); errValidation != nil {
		if err := svc.repo.CallUpdateJob(jobId, &updateJobRequest{
			Status:       helper.StringPtr(constant.JobStatusDone),
			SuccessCount: helper.IntPtr(1),
			EndAt:        helper.TimePtr(time.Now()),
		}); err != nil {
			return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveJob)
		}

		if err := svc.repo.CallUpdateJobDetail(jobId, jobDetailId, &updateJobDetailRequest{
			Message:    helper.StringPtr(errValidation.Error()),
			InProgress: helper.BoolPtr(false),
			Status:     helper.StringPtr(constant.JobStatusFailed),
		}); err != nil {
			return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveDetail)
		}

		return apperror.BadRequest(errValidation.Error())
	}

	return nil
}

func (svc *service) processPhoneLiveStatus(apiKey, jobId, jobDetailId string, jobDetail *mstPhoneLiveStatusJobDetail) (*model.ProCatAPIResponse[phoneLiveStatusRespData], error) {
	req := &phoneLiveStatusRequest{
		PhoneNumber: jobDetail.PhoneNumber,
		TrxId:       strconv.FormatUint(uint64(jobDetail.JobId), 10),
	}

	result, err := svc.repo.CallPhoneLiveStatusAPI(apiKey, req)
	if err != nil {
		if err := svc.repo.CallUpdateJobDetail(jobId, jobDetailId, &updateJobDetailRequest{
			Message:    helper.StringPtr(err.Error()),
			InProgress: helper.BoolPtr(false),
			Status:     helper.StringPtr(constant.JobStatusError),
		}); err != nil {
			return nil, apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveDetail)
		}

		return nil, apperror.MapRepoError(err, "failed to process phone live status request")
	}

	return result, nil
}

func (svc *service) processAndUpdatePhoneLiveStatus(apiKey, jobId, jobDetailId string, detail *mstPhoneLiveStatusJobDetail) error {
	if err := validator.ValidateStruct(detail); err != nil {
		_ = svc.repo.CallUpdateJobDetail(jobId, jobDetailId, &updateJobDetailRequest{
			Message:    helper.StringPtr(err.Error()),
			InProgress: helper.BoolPtr(false),
			Status:     helper.StringPtr(constant.JobStatusError),
		})

		return apperror.BadRequest(err.Error())
	}

	resp, err := svc.processPhoneLiveStatus(apiKey, jobId, jobDetailId, detail)
	if err != nil {
		return err
	}

	if err := svc.updateProcessedDetail(jobId, jobDetailId, resp); err != nil {
		return err
	}

	return nil
}

func parseLiveStatusData(liveStatus string) (subscriberStatus, deviceStatus string, err error) {
	if liveStatus == "" {
		return "", "", nil
	}

	parts := strings.Split(liveStatus, ",")

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func (svc *service) updateProcessedDetail(jobId, jobDetailId string, phoneLiveResp *model.ProCatAPIResponse[phoneLiveStatusRespData]) error {
	status := "success"
	if len(phoneLiveResp.Data.Errors) > 0 && phoneLiveResp.Data.Errors[0].Code == -60001 {
		status = "fail"
	}

	subcriberStatus, deviceStatus, err := parseLiveStatusData(phoneLiveResp.Data.LiveStatus)
	if err != nil {
		return apperror.Internal("failed to parse live status data", err)
	}

	if err := svc.repo.CallUpdateJobDetail(jobId, jobDetailId, &updateJobDetailRequest{
		Status:           helper.StringPtr(status),
		InProgress:       helper.BoolPtr(false),
		SubscriberStatus: helper.StringPtr(subcriberStatus),
		DeviceStatus:     helper.StringPtr(deviceStatus),
		PhoneType:        helper.StringPtr(phoneLiveResp.Data.PhoneType),
		Operator:         helper.StringPtr(phoneLiveResp.Data.Operator),
		PricingStrategy:  helper.StringPtr(phoneLiveResp.PricingStrategy),
		TransactionId:    helper.StringPtr(phoneLiveResp.TransactionId),
	}); err != nil {
		return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveDetail)
	}

	return nil
}

func (svc *service) finalizeJob(jobId string) error {
	count, err := svc.repo.CallGetProcessedCountAPI(jobId)
	if err != nil {
		return apperror.MapRepoError(err, "failed to get processed count request")
	}

	if err := svc.repo.CallUpdateJob(jobId, &updateJobRequest{
		Status:       helper.StringPtr(constant.JobStatusDone),
		SuccessCount: helper.IntPtr(int(count.SuccessCount)),
		EndAt:        helper.TimePtr(time.Now()),
	}); err != nil {
		return apperror.MapRepoError(err, constant.ErrMsgUpdatePhoneLiveJob)
	}

	return nil
}
