package livestatus

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"front-office/helper"
	"log"
	"strconv"
	"sync"
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
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

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

	jobID, err := ctrl.Svc.CreateJob(liveStatusRequests, userID, companyID, totalData)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	jobDetails, err := ctrl.Svc.GetJobDetailsByJobID(jobID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(jobDetails))

	var successRequestTotal int
	var countJob int
	countJob = 0
	for _, jobDetail := range jobDetails {
		countJob += 1
		if countJob == 100 {
			time.Sleep(1 * time.Second)
			countJob = 0
		}

		wg.Add(1)
		go func(jobDetail *JobDetail) {
			defer wg.Done()

			var err error
			if errValid := validator.ValidateStruct(jobDetail); errValid != nil {
				err = ctrl.Svc.UpdateInvalidJobDetail(jobDetail.ID, errValid.Error())
				if err != nil {
					errChan <- err
					return
				}
			} else {
				err = ctrl.Svc.ProcessJobDetails(jobDetail)
				if err != nil {
					errChan <- err
					return
				}
			}

			// update count success pada tabel job
			successRequestTotal, err = ctrl.Svc.CountOnProcessJobDetails(jobDetail.JobID, false)
			if err != nil {
				errChan <- err
				return
			}

			updateReq := UpdateJobRequest{
				Total: &successRequestTotal,
			}
			err = ctrl.Svc.UpdateJob(jobDetail.JobID, &updateReq)
			if err != nil {
				errChan <- err
				return
			}
		}(jobDetail)
	}

	// jika tidak ada job details dengan status 'error', update status pada job menjadi 'done'
	// failedJobDetails, err := ctrl.Svc.GetFailedJobDetails(jobID)
	// if err != nil {
	// 	log.Println("Error GetFailedJobDetails : ", err.Error())
	// }

	// if failedJobDetails != nil && len(failedJobDetails) == 0 {
	// 	doneStatus := "done"
	// 	now := time.Now()

	// 	updateReq := UpdateJobRequest{
	// 		Status: &doneStatus,
	// 		EndAt:  &now,
	// 	}

	// 	err := ctrl.Svc.UpdateJob(jobID, &updateReq)
	// 	if err != nil {
	// 		statusCode, resp := helper.GetError(err.Error())
	// 		return c.Status(statusCode).JSON(resp)
	// 	}
	// }

	wg.Wait()
	close(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	default:
		{
			fmt.Println("No errors found in job processing")
		}
	}

	doneStatus := "done"
	now := time.Now()

	updateReq := UpdateJobRequest{
		Status: &doneStatus,
		EndAt:  &now,
	}

	err = ctrl.Svc.UpdateJob(jobID, &updateReq)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	// todo: jika semua request sukses, hapus job pada temp tabel
	// err = ctrl.Svc.DeleteJob(jobID)
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
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	jobs, err := ctrl.Svc.GetJobs(page, size, userID, companyID, startDate, endDate, uint(tierLevel))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetJobsTotal(userID, companyID, startDate, endDate, uint(tierLevel))

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
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)

	totalData, _ := ctrl.Svc.GetJobsTotalByRangeDate(userID, companyID, startDate, endDate, uint(tierLevel))
	totalSubscriberActive, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userID, companyID, startDate, endDate, "subscriber_status", "ACTIVE", uint(tierLevel))
	totalDeviceReachable, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userID, companyID, startDate, endDate, "device_status", "REACHABLE", uint(tierLevel))
	totalMobilePhone, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userID, companyID, startDate, endDate, "data", "MOBILE", uint(tierLevel))
	totalFixedLine, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(userID, companyID, startDate, endDate, "data", "FIXED_LINE", uint(tierLevel))
	totalDataPercentageSuccess, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(userID, companyID, startDate, endDate, "success", uint(tierLevel))
	totalDataPercentageFail, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(userID, companyID, startDate, endDate, "fail", uint(tierLevel))
	totalDataPercentageError, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(userID, companyID, startDate, endDate, "error", uint(tierLevel))

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
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)

	jobDetails, err := ctrl.Svc.GetJobDetailsByRangeDate(userID, companyID, startDate, endDate, uint(tierLevel))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	header := []string{"Phone Number", "Phone Type", "Operator", "Device Status", "Subscriber Status", "Message"}
	if err := w.Write(header); err != nil {
		statusCode, resp := helper.GetError("Failed to write CSV header")
		return c.Status(statusCode).JSON(resp)
	}

	for _, record := range jobDetails {
		row := []string{record.PhoneNumber, record.PhoneType, record.Operator, record.DeviceStatus, record.SubscriberStatus, record.Message}
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
	jobID := c.Params("id")
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	jobIDUint, _ := strconv.ParseUint(jobID, 10, 32)

	_, err := ctrl.Svc.GetJobByIDAndUserID(uint(jobIDUint), uint(tierLevel), userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	jobs, err := ctrl.Svc.GetJobDetailsWithPagination(page, size, keyword, uint(jobIDUint))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetJobDetailsWithPaginationTotal(keyword, uint(jobIDUint))
	subscriberStatuscon, _ := ctrl.Svc.GetJobDetailsPercentage("subscriber_status", "ACTIVE", uint(jobIDUint))
	deviceStatusReach, _ := ctrl.Svc.GetJobDetailsPercentage("device_status", "REACHABLE", uint(jobIDUint))
	totalDataPercentage, _ := ctrl.Svc.GetJobDetailsWithPaginationTotalPercentage(uint(jobIDUint), "success")
	totalDataPercentageFail, _ := ctrl.Svc.GetJobDetailsWithPaginationTotalPercentage(uint(jobIDUint), "fail")
	totalDataPercentageError, _ := ctrl.Svc.GetJobDetailsWithPaginationTotalPercentage(uint(jobIDUint), "error")

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
	jobIDs, err := ctrl.Svc.GetJobWithIncompleteStatus()
	if err != nil {
		log.Println("Error GetJobWithIncompleStatus : ", err.Error())
		return
	}

	if len(jobIDs) == 0 {
		log.Println("No incomplete job found")
		return
	}

	for _, jobID := range jobIDs {
		jobDetails, err := ctrl.Svc.GetFailedJobDetails(jobID)
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
			successRequestTotal, err := ctrl.Svc.CountOnProcessJobDetails(jobDetail.JobID, false)
			if err != nil {
				log.Println("Error CountOnProcessJobDetails : ", err.Error())
			}

			updateReq := UpdateJobRequest{
				Total: &successRequestTotal,
			}
			err = ctrl.Svc.UpdateJob(jobDetail.JobID, &updateReq)
			if err != nil {
				log.Println("Error UpdateJob : ", err.Error())
			}
		}

		// jika tidak ada lagi job detail yang diproses dalam sebuah job, update status pada job menjadi 'done'
		jobDetailIDs, err := ctrl.Svc.GetOnProcessJobDetails(jobID, true)
		if err != nil {
			log.Println("Error GetOnProcessJobDetails : ", err.Error())
		}

		if len(jobDetailIDs) == 0 {
			doneStatus := "done"
			now := time.Now()

			updateReq := UpdateJobRequest{
				Status: &doneStatus,
				EndAt:  &now,
			}

			err := ctrl.Svc.UpdateJob(jobID, &updateReq)
			if err != nil {
				log.Println("Error UpdateJob : ", err.Error())
			}
		}
	}

	// todo: jika semua request sukses, hapus job pada temp tabel
	// err = ctrl.Svc.DeleteJob(jobID)
	// if err != nil {
	// 	log.Println("Error DeleteJob : ", err.Error())
	// }
}

func (ctrl *controller) GetJobDetailsExport(c *fiber.Ctx) error {
	jobID := c.Params("id")
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	jobIDUint, _ := strconv.ParseUint(jobID, 10, 32)

	// Get Job by ID
	_, err := ctrl.Svc.GetJobByIDAndUserID(uint(jobIDUint), uint(tierLevel), userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	// Get JobDetails By JobID
	jobDetails, err := ctrl.Svc.GetJobDetailsByJobIDExport(uint(jobIDUint))
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

	var filename string = fmt.Sprintf("jobs_detail_%s.csv", jobID)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendStream(bytes.NewReader(buf.Bytes()))
}
