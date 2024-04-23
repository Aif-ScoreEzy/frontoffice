package livestatus

import (
	"fmt"
	"front-office/helper"
	"log"
	"strconv"

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
	ReprocessUnsuccessfulJobDetails()
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
	for _, jobDetail := range jobDetails {
		successRequestTotal, err = ctrl.Svc.ProcessJobDetails(jobDetail, successRequestTotal)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
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

	jobs, err := ctrl.Svc.GetJobs(page, size)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetJobsTotal()

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

	fullResponsePage := map[string]interface{}{
		"total_data":    totalData,
		"subs_active":   subscriberStatuscon,
		"dev_reachable": deviceStatusReach,
		"data":          jobs,
	}

	resp := helper.ResponseSuccess(
		"success",
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) ReprocessUnsuccessfulJobDetails() {
	jobDetails, err := ctrl.Svc.GetUnprocessedJobDetails()
	if err != nil {
		log.Println("Error GetUnprocessedJobDetails : ", err.Error())
	}

	if jobDetails == nil {
		log.Println("No unprocessed job details found")
		return
	}

	var successRequestTotal int
	for _, jobDetail := range jobDetails {
		if err := ctrl.Svc.UpdateProcessedJobDetail(jobDetail.ID); err != nil {
			log.Println("Error UpdateProcessedJobDetail : ", err.Error())
		}

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
