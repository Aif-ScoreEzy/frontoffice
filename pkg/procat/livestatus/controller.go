package livestatus

import (
	"fmt"
	"front-office/app/config"
	"front-office/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service, cfg *config.Config) Controller {
	return &controller{Svc: service, Cfg: cfg}
}

type controller struct {
	Cfg *config.Config
	Svc Service
}

type Controller interface {
	BulkSearch(c *fiber.Ctx) error
	GetJobs(c *fiber.Ctx) error
	GetJobDetails(c *fiber.Ctx) error
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
	apiKey := ctrl.Cfg.Env.ApiKeyLiveStatus

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

	// batchSize := 10
	// liveStatusResponses, err := ctrl.Svc.ProcessJobDetails(apiKey, jobID, jobDetails, batchSize)
	// if err != nil {
	// 	statusCode, resp := helper.GetError(err.Error())
	// 	return c.Status(statusCode).JSON(resp)
	// }

	var successRequestTotal int
	var liveStatusResponse *LiveStatusResponse
	jobIDStr := strconv.FormatUint(uint64(jobID), 10)
	for _, jobDetail := range jobDetails {
		request := &LiveStatusRequest{
			PhoneNumber: jobDetail.PhoneNumber,
			TrxID:       jobIDStr,
		}

		liveStatusResponse, err = ctrl.Svc.CreateLiveStatus(request, apiKey)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		// todo: jika status code 200 kirim job detail ke aifcore

		dataMap := liveStatusResponse.Data.(map[string]interface{})
		dataLiveMap := dataMap["live"].(map[string]interface{})
		subscriberStatus := fmt.Sprintf("%v", dataLiveMap["subscriber_status"])
		deviceStatus := fmt.Sprintf("%v", dataLiveMap["device_status"])

		// todo: jika status code 200 maka hapus job detail pada temp tabel. Sampai aifcore menyediakan API untuk get job details, untuk sementara jika status code 200 lakukan update subcriber_status dan device_status pada job detail
		if liveStatusResponse.StatusCode == 200 {
			successRequestTotal += 1
			// err = ctrl.Svc.DeleteJobDetail(jobDetail.ID)
			err = ctrl.Svc.UpdateSucceededJobDetail(jobDetail.ID, subscriberStatus, deviceStatus)
			if err != nil {

				statusCode, resp := helper.GetError(err.Error())
				return c.Status(statusCode).JSON(resp)
			}
		} else {
			err = ctrl.Svc.UpdateFailedJobDetail(jobID, jobDetail.Sequence)
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

	// todo: jika dari aifcore sudah tersedia api untuk get jobs, hapus program update job
	err = ctrl.Svc.UpdateJob(jobID, successRequestTotal)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

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
