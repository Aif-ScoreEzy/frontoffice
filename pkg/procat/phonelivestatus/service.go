package phonelivestatus

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/member"
	"log"
	"mime/multipart"
	"strconv"
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
	CreateJob(memberId, companyId string, request *createJobRequest) (*createJobResponseData, error)
	GetPhoneLiveStatusJob(filter *phoneLiveStatusFilter) (*jobListRespData, error)
	GetAllPhoneLiveStatusDetails(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	GetPhoneLiveStatusDetailsByRangeDate(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error)
	GetJobsSummary(filter *phoneLiveStatusFilter) (*jobsSummaryRespData, error)
	GetPhoneLiveStatusDetailsSummary(filter *phoneLiveStatusFilter) (*jobDetailRespData, error)
	ExportJobsSummary(filter *phoneLiveStatusFilter, buf *bytes.Buffer) (string, error)
	ExportJobDetails(filter *phoneLiveStatusFilter, buf *bytes.Buffer) (string, error)
	UpdateJob(jobId uint, req *updateJobRequest) (*model.AifcoreAPIResponse[any], error)
	UpdateJobDetail(jobId, jobDetailId uint, req *updateJobDetailRequest) (*model.AifcoreAPIResponse[any], error)
	ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error
	BulkProcessPhoneLiveStatus(memberId, companyId string, fileHeader *multipart.FileHeader) error
}

func (svc *service) CreateJob(memberId, companyId string, request *createJobRequest) (*createJobResponseData, error) {
	job, err := svc.repo.CallCreateJobAPI(memberId, companyId, request)
	if err != nil {
		return nil, apperror.MapRepoError(err, "create job failed")
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
		return nil, apperror.MapRepoError(err, "failed to fetch phone live status job detail")
	}

	return jobDetails, nil
}

func (svc *service) GetAllPhoneLiveStatusDetails(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error) {
	jobDetails, err := svc.repo.CallGetAllJobDetailsAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch phone live status job detail")
	}

	return jobDetails, nil
}

func (svc *service) GetPhoneLiveStatusDetailsByRangeDate(filter *phoneLiveStatusFilter) ([]*mstPhoneLiveStatusJobDetail, error) {
	jobDetail, err := svc.repo.CallGetJobDetailsByRangeDateAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch phone live status job detail")
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

func (svc *service) UpdateJob(jobId uint, req *updateJobRequest) (*model.AifcoreAPIResponse[any], error) {
	jobIdStr := strconv.FormatUint(uint64(jobId), 10)
	response, err := svc.repo.CallUpdateJob(jobIdStr, req)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}

func (svc *service) UpdateJobDetail(jobId, jobDetailId uint, req *updateJobDetailRequest) (*model.AifcoreAPIResponse[any], error) {
	jobIdStr := strconv.FormatUint(uint64(jobId), 10)
	jobDetailIdStr := strconv.FormatUint(uint64(jobId), 10)
	response, err := svc.repo.CallUpdateJobDetail(jobIdStr, jobDetailIdStr, req)
	if err != nil {
		return nil, err
	}

	return helper.ParseAifcoreAPIResponse[any](response)
}

func (svc *service) ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error {
	// member, err := svc.memberRepo.CallGetMemberAPI(&member.FindUserQuery{
	// 	Id:        memberId,
	// 	CompanyId: companyId,
	// })
	// if err != nil {
	// 	return apperror.MapRepoError(err, "failed to fetch ")
	// }

	response, err := svc.repo.CallPhoneLiveStatusAPI(memberId, companyId, req)
	log.Println("phone live status resss==> ", response, err)

	if err != nil {
		return err
	}

	return nil
}

func (svc *service) BulkProcessPhoneLiveStatus(memberId, companyId string, fileHeader *multipart.FileHeader) error {
	_, err := svc.repo.CallBulkPhoneLiveStatusAPI(memberId, companyId, fileHeader)
	if err != nil {
		return err
	}

	return nil
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
