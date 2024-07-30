package livestatus

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"front-office/helper"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

func NewController(service Service) Controller {
	return &controller{Svc: service}
}

type controller struct {
	Svc Service
}

type Controller interface {
	BulkSearch(c *fiber.Ctx) error
	GetJobs(c *fiber.Ctx) error
	GetJobDetails(c *fiber.Ctx) error
	GetJobsSummary(c *fiber.Ctx) error
	ExportJobsSummary(c *fiber.Ctx) error
	ReprocessFailedJobDetails()
	GetJobDetailsExport(c *fiber.Ctx) error
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	file, err := c.FormFile("file")
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	expectedHeaders := []string{"phone_number"}
	csvData, totalData, err := helper.ParseCSVFile(file, expectedHeaders)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var liveStatusRequest LiveStatusRequest
	var liveStatusRequests []LiveStatusRequest
	for i, request := range csvData {
		if i == 0 {
			continue
		}

		liveStatusRequest.PhoneNumber = request[0]

		liveStatusRequests = append(liveStatusRequests, liveStatusRequest)
	}

	jobId, err := ctrl.Svc.CreateJob(liveStatusRequests, userId, companyId, totalData)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	jobDetails, err := ctrl.Svc.GetJobDetailsByJobId(jobId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var successRequestTotal int
	for _, jobDetail := range jobDetails {
		if errValid := validator.ValidateStruct(jobDetail); errValid != nil {
			err = ctrl.Svc.UpdateInvalidJobDetail(jobDetail.Id, errValid.Error())
			if err != nil {
				statusCode, resp := helper.GetError(err.Error())
				return c.Status(statusCode).JSON(resp)
			}
		} else {
			err = ctrl.Svc.ProcessJobDetails(jobDetail)
			if err != nil {
				statusCode, resp := helper.GetError(err.Error())
				return c.Status(statusCode).JSON(resp)
			}
		}

		// update count success pada tabel job
		successRequestTotal, err = ctrl.Svc.CountOnProcessJobDetails(jobDetail.JobId, false)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		updateReq := UpdateJobRequest{
			Total: &successRequestTotal,
		}
		err = ctrl.Svc.UpdateJob(jobDetail.JobId, &updateReq)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	// jika tidak ada job details dengan status 'error', update status pada job menjadi 'done'
	failedJobDetails, err := ctrl.Svc.GetFailedJobDetails(jobId)
	if err != nil {
		log.Println("Error GetFailedJobDetails : ", err.Error())
	}

	if failedJobDetails != nil && len(failedJobDetails) == 0 {
		doneStatus := "done"
		now := time.Now()

		updateReq := UpdateJobRequest{
			Status: &doneStatus,
			EndAt:  &now,
		}

		err := ctrl.Svc.UpdateJob(jobId, &updateReq)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	// todo: jika semua request sukses, hapus job pada temp tabel
	// err = ctrl.Svc.DeleteJob(jobId)
	// if err != nil {
	// 	statusCode, resp := helper.GetError(err.Error())
	// 	return c.Status(statusCode).JSON(resp)
	// }

	dataResponse := ResponseSuccess{
		Success:   successRequestTotal,
		TotalData: totalData,
	}

	resp := helper.ResponseSuccess(
		"success",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetJobs(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	jobs, err := ctrl.Svc.GetJobs(page, size, userId, companyId, startDate, endDate, uint(tierLevel))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetJobsTotal(userId, companyId, startDate, endDate, uint(tierLevel))

	data := GetJobsResponse{
		TotalData: totalData,
		Jobs:      jobs,
	}

	resp := helper.ResponseSuccess(
		"success to get jobs",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetJobsSummary(c *fiber.Ctx) error {
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)

	totalData, _ := ctrl.Svc.GetJobsTotalByRangeDate(userId, companyId, startDate, endDate, uint(tierLevel))
	totalSubscriberActive, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startDate, endDate, "subscriber_status", "ACTIVE", uint(tierLevel))
	totalDeviceReachable, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startDate, endDate, "device_status", "REACHABLE", uint(tierLevel))
	totalMobilePhone, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startDate, endDate, "data", "MOBILE", uint(tierLevel))
	totalFixedLine, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userId, companyId, startDate, endDate, "data", "FIXED_LINE", uint(tierLevel))
	totalDataPercentageSuccess, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(userId, companyId, startDate, endDate, "success", uint(tierLevel))
	totalDataPercentageFail, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(userId, companyId, startDate, endDate, "fail", uint(tierLevel))
	totalDataPercentageError, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(userId, companyId, startDate, endDate, "error", uint(tierLevel))

	data := JobSummaryResponse{
		TotalData:        totalData,
		TotalDataSuccess: totalDataPercentageSuccess,
		TotalDataFail:    totalDataPercentageFail,
		TotalDataError:   totalDataPercentageError,
		SubscriberActive: totalSubscriberActive,
		DeviceReachable:  totalDeviceReachable,
		Mobile:           totalMobilePhone,
		FixedLine:        totalFixedLine,
	}

	resp := helper.ResponseSuccess(
		"success",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) ExportJobsSummary(c *fiber.Ctx) error {
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)

	jobDetails, err := ctrl.Svc.GetJobDetailsByRangeDate(userId, companyId, startDate, endDate, uint(tierLevel))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	header := []string{"Phone Number", "Phone Type", "Operator", "Device Status", "Subscriber Status"}
	if err := w.Write(header); err != nil {
		statusCode, resp := helper.GetError("Failed to write CSV header")
		return c.Status(statusCode).JSON(resp)
	}

	for _, record := range jobDetails {
		row := []string{record.PhoneNumber, record.PhoneType, record.Operator, record.DeviceStatus, record.SubscriberStatus}
		if err := w.Write(row); err != nil {
			statusCode, resp := helper.GetError("Failed to write CSV data")
			return c.Status(statusCode).JSON(resp)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		statusCode, resp := helper.GetError("Failed to flush CSV data")
		return c.Status(statusCode).JSON(resp)
	}

	var filename string
	if endDate != "" && endDate != startDate {
		filename = fmt.Sprintf("jobs_summary_%s_until_%s.csv", startDate, endDate)
	} else {
		filename = fmt.Sprintf("jobs_summary_%s.csv", startDate)
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}

func (ctrl *controller) GetJobDetails(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	keyword := c.Query("keyword", "")
	jobId := c.Params("id")
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	jobIdUint, _ := strconv.ParseUint(jobId, 10, 32)

	_, err := ctrl.Svc.GetJobByIdAndUserId(uint(jobIdUint), uint(tierLevel), userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	jobs, err := ctrl.Svc.GetJobDetailsWithPagination(page, size, keyword, uint(jobIdUint))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetJobDetailsWithPaginationTotal(keyword, uint(jobIdUint))
	subscriberStatuscon, _ := ctrl.Svc.GetJobDetailsPercentage("subscriber_status", "ACTIVE", uint(jobIdUint))
	deviceStatusReach, _ := ctrl.Svc.GetJobDetailsPercentage("device_status", "REACHABLE", uint(jobIdUint))
	totalDataPercentage, _ := ctrl.Svc.GetJobDetailsWithPaginationTotalPercentage(uint(jobIdUint), "success")
	totalDataPercentageFail, _ := ctrl.Svc.GetJobDetailsWithPaginationTotalPercentage(uint(jobIdUint), "fail")
	totalDataPercentageError, _ := ctrl.Svc.GetJobDetailsWithPaginationTotalPercentage(uint(jobIdUint), "error")

	dataResponse := JobDetailResponse{
		TotalData:        totalData,
		TotalDataSuccess: totalDataPercentage,
		TotalDataFail:    totalDataPercentageFail,
		TotalDataError:   totalDataPercentageError,
		SubscriberActive: subscriberStatuscon,
		DeviceReachable:  deviceStatusReach,
		JobDetails:       jobs,
	}

	resp := helper.ResponseSuccess(
		"success to get job details",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) ReprocessFailedJobDetails() {
	jobIds, err := ctrl.Svc.GetJobWithIncompleteStatus()
	if err != nil {
		log.Println("Error GetJobWithIncompleStatus : ", err.Error())
		return
	}

	if len(jobIds) == 0 {
		log.Println("No incomplete job found")
		return
	}

	for _, jobId := range jobIds {
		jobDetails, err := ctrl.Svc.GetFailedJobDetails(jobId)
		if err != nil {
			log.Println("Error GetFailedJobDetails : ", err.Error())
		}

		if jobDetails != nil && len(jobDetails) == 0 {
			log.Println("No failed job details found")
		}

		for _, jobDetail := range jobDetails {
			err = ctrl.Svc.ProcessJobDetails(jobDetail)
			if err != nil {
				log.Println("Error ProcessJobDetails : ", err.Error())
			}

			// update count success pada tabel job
			successRequestTotal, err := ctrl.Svc.CountOnProcessJobDetails(jobDetail.JobId, false)
			if err != nil {
				log.Println("Error CountOnProcessJobDetails : ", err.Error())
			}

			updateReq := UpdateJobRequest{
				Total: &successRequestTotal,
			}
			err = ctrl.Svc.UpdateJob(jobDetail.JobId, &updateReq)
			if err != nil {
				log.Println("Error UpdateJob : ", err.Error())
			}
		}

		// jika tidak ada lagi job detail yang diproses dalam sebuah job, update status pada job menjadi 'done'
		jobDetailIds, err := ctrl.Svc.GetOnProcessJobDetails(jobId, true)
		if err != nil {
			log.Println("Error GetOnProcessJobDetails : ", err.Error())
		}

		if len(jobDetailIds) == 0 {
			doneStatus := "done"
			now := time.Now()

			updateReq := UpdateJobRequest{
				Status: &doneStatus,
				EndAt:  &now,
			}

			err := ctrl.Svc.UpdateJob(jobId, &updateReq)
			if err != nil {
				log.Println("Error UpdateJob : ", err.Error())
			}
		}
	}

	// todo: jika semua request sukses, hapus job pada temp tabel
	// err = ctrl.Svc.DeleteJob(jobId)
	// if err != nil {
	// 	log.Println("Error DeleteJob : ", err.Error())
	// }
}

func (ctrl *controller) GetJobDetailsExport(c *fiber.Ctx) error {
	jobId := c.Params("id")
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	jobIdUint, _ := strconv.ParseUint(jobId, 10, 32)

	// Get Job by Id
	_, err := ctrl.Svc.GetJobByIdAndUserId(uint(jobIdUint), uint(tierLevel), userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	// Get JobDetails By JobId
	jobDetails, err := ctrl.Svc.GetJobDetailsByJobIdExport(uint(jobIdUint))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	header := []string{"Phone Number", "Subscriber Status", "Device Status", "Status", "Operator", "Phone Type"}
	if err := w.Write(header); err != nil {
		statusCode, resp := helper.GetError("Failed to write CSV header")
		return c.Status(statusCode).JSON(resp)
	}

	for _, record := range jobDetails {
		row := []string{record.PhoneNumber, record.SubscriberStatus, record.DeviceStatus, record.Status, record.Operator, record.PhoneType}
		if err := w.Write(row); err != nil {
			statusCode, resp := helper.GetError("Failed to write CSV data")
			return c.Status(statusCode).JSON(resp)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		statusCode, resp := helper.GetError("Failed to flush CSV data")
		return c.Status(statusCode).JSON(resp)
	}

	var filename string = fmt.Sprintf("jobs_detail_%s.csv", jobId)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}
