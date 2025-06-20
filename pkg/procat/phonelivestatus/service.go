package phonelivestatus

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror/mapper"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

type service struct {
	repo Repository
}

type Service interface {
	CreateJob(memberId, companyId string, request *createJobRequest) (*createJobResponseData, error)
	GetPhoneLiveStatusJob(filter *PhoneLiveStatusFilter) (*APIResponse[JobListResponse], error)
	GetAllPhoneLiveStatusDetails(filter *PhoneLiveStatusFilter) (*APIResponse[[]MstPhoneLiveStatusJobDetail], error)
	GetPhoneLiveStatusDetailsByRangeDate(filter *PhoneLiveStatusFilter) (*APIResponse[[]MstPhoneLiveStatusJobDetail], error)
	GetJobsSummary(filter *PhoneLiveStatusFilter) (*APIResponse[JobsSummaryResponse], error)
	GetPhoneLiveStatusDetailsSummary(filter *PhoneLiveStatusFilter) (*APIResponse[JobDetailsResponse], error)
	ExportJobsSummary(data []MstPhoneLiveStatusJobDetail, filter *PhoneLiveStatusFilter, buf *bytes.Buffer) (string, error)
	UpdateJob(jobId uint, req *updateJobRequest) (*model.AifcoreAPIResponse[any], error)
	UpdateJobDetail(jobId, jobDetailId uint, req *updateJobDetailRequest) (*model.AifcoreAPIResponse[any], error)
	ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error
	BulkProcessPhoneLiveStatus(memberId, companyId string, fileHeader *multipart.FileHeader) error
}

func (svc *service) CreateJob(memberId, companyId string, request *createJobRequest) (*createJobResponseData, error) {
	res, err := svc.repo.CallCreateJobAPI(memberId, companyId, request)
	if err != nil {
		return nil, mapper.MapRepoError(err, "create job failed")
	}

	return &res.Data, err
}

func (svc *service) GetPhoneLiveStatusJob(filter *PhoneLiveStatusFilter) (*APIResponse[JobListResponse], error) {
	response, err := svc.repo.CallGetPhoneLiveStatusJobAPI(filter)
	if err != nil {
		return nil, err
	}

	return parseGenericResponse[JobListResponse](response)
}

func (svc *service) GetPhoneLiveStatusDetailsSummary(filter *PhoneLiveStatusFilter) (*APIResponse[JobDetailsResponse], error) {
	response, err := svc.repo.CallGetJobDetailsAPI(filter)
	if err != nil {
		return nil, err
	}

	return parseGenericResponse[JobDetailsResponse](response)
}

func (svc *service) GetAllPhoneLiveStatusDetails(filter *PhoneLiveStatusFilter) (*APIResponse[[]MstPhoneLiveStatusJobDetail], error) {
	response, err := svc.repo.CallGetAllJobDetailsAPI(filter)
	if err != nil {
		return nil, err
	}

	return parseGenericResponse[[]MstPhoneLiveStatusJobDetail](response)
}

func (svc *service) GetPhoneLiveStatusDetailsByRangeDate(filter *PhoneLiveStatusFilter) (*APIResponse[[]MstPhoneLiveStatusJobDetail], error) {
	response, err := svc.repo.CallGetJobDetailsByRangeDateAPI(filter)
	if err != nil {
		return nil, err
	}

	return parseGenericResponse[[]MstPhoneLiveStatusJobDetail](response)
}

func (svc *service) GetJobsSummary(filter *PhoneLiveStatusFilter) (*APIResponse[JobsSummaryResponse], error) {
	response, err := svc.repo.CallGetJobsSummary(filter)
	if err != nil {
		return nil, err
	}

	return parseGenericResponse[JobsSummaryResponse](response)
}

func (svc *service) ExportJobsSummary(data []MstPhoneLiveStatusJobDetail, filter *PhoneLiveStatusFilter, buf *bytes.Buffer) (string, error) {
	w := csv.NewWriter(buf)

	header := []string{"Phone Number", "Subscriber Status", "Device Status", "Status", "Operator", "Phone Type"}
	if err := w.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header")
	}

	for _, record := range data {
		row := []string{record.PhoneNumber, record.SubscriberStatus, record.DeviceStatus, record.Status, record.Operator, record.PhoneType}
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV data")
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV data")
	}

	var filename string
	if filter.EndDate != "" && filter.EndDate != filter.StartDate {
		filename = fmt.Sprintf("jobs_summary_%s_until_%s.csv", filter.StartDate, filter.EndDate)
	} else {
		filename = fmt.Sprintf("job_summary_%s.csv", filter.StartDate)
	}

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

func parseGenericResponse[T any](response *http.Response) (*APIResponse[T], error) {
	var apiResponse APIResponse[T]

	if response == nil {
		return nil, errors.New("nil response")
	}

	dataBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, &apiResponse); err != nil {
		return nil, err
	}

	apiResponse.StatusCode = response.StatusCode

	return &apiResponse, nil
}
