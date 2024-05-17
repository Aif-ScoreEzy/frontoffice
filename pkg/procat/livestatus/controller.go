package livestatus

import (
	"fmt"
	"front-office/helper"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
	ReprocessFailedJobDetails()
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
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

	jobID, err := ctrl.Svc.CreateJob(liveStatusRequests, totalData)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	jobDetails, err := ctrl.Svc.GetJobDetails(jobID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var successRequestTotal int
	for i, jobDetail := range jobDetails {
		successRequestTotal, err = ctrl.Svc.ProcessJobDetails(jobDetail, successRequestTotal)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if i == totalData-1 {
			doneStatus := "done"
			now := time.Now()

			updateReq := UpdateJobRequest{
				Status: &doneStatus,
				EndAt:  &now,
			}

			err := ctrl.Svc.UpdateJob(jobID, &updateReq)
			if err != nil {
				statusCode, resp := helper.GetError(err.Error())
				return c.Status(statusCode).JSON(resp)
			}
		}
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
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	jobs, err := ctrl.Svc.GetJobs(page, size, startDate, endDate)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetJobsTotal(startDate, endDate)

	fullResponsePage := map[string]interface{}{
		"total_data": totalData,
		"data":       jobs,
	}

	resp := helper.ResponseSuccess(
		"success",
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetJobsSummary(c *fiber.Ctx) error {
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	totalData, _ := ctrl.Svc.GetJobsTotalByRangeDate(startDate, endDate)
	totalSubscriberActive, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(startDate, endDate, "subscriber_status", "ACTIVE")
	totalDeviceReachable, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(startDate, endDate, "device_status", "REACHABLE")
	totalMobilePhone, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(startDate, endDate, "data", "MOBILE")
	totalFixedLine, _ := ctrl.Svc.GetJobDetailsPercentageByDataAndRangeDate(startDate, endDate, "data", "FIXED_LINE")
	totalDataPercentageSuccess, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(startDate, endDate, "success")
	totalDataPercentageFail, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(startDate, endDate, "fail")
	totalDataPercentageError, _ := ctrl.Svc.GetJobDetailsTotalPercentageByRangeDate(startDate, endDate, "error")

	fullResponsePage := JobSummaryResponse{
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
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetJobDetails(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	keyword := c.Query("keyword", "")
	jobID := c.Params("id")

	jobIDUint, _ := strconv.ParseUint(fmt.Sprintf("%v", jobID), 10, 64)
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
	jobDetails, err := ctrl.Svc.GetFailedJobDetails()
	if err != nil {
		log.Println("Error GetFailedJobDetails : ", err.Error())
	}

	if jobDetails == nil {
		log.Println("No failed job details found")
		return
	}

	var successRequestTotal int
	for _, jobDetail := range jobDetails {
		job, _ := ctrl.Svc.GetJobByID(jobDetail.JobID)
		successRequestTotal = job.Success

		_, _ = ctrl.Svc.ProcessJobDetails(jobDetail, successRequestTotal)
		if err != nil {
			log.Println("Error ProcessJobDetails : ", err.Error())
		}
	}

	// todo: jika semua request sukses, hapus job pada temp tabel
	// err = ctrl.Svc.DeleteJob(jobID)
	// if err != nil {
	// 	log.Println("Error DeleteJob : ", err.Error())
	// }
}
