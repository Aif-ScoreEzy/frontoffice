package livestatus

import (
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
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
	// apiKey := c.Get("X-AIF-KEY")
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

		// todo: jika sukses kirim ke aifcore

		dataMap := liveStatusResponse.Data.(map[string]interface{})
		errors := dataMap["errors"].([]interface{})
		if len(errors) == 0 {
			successRequestTotal += 1
			err = ctrl.Svc.DeleteJobDetail(jobDetail.ID)
			if err != nil {
				statusCode, resp := helper.GetError(err.Error())
				return c.Status(statusCode).JSON(resp)
			}
		}
	}

	// todo: jika semua request sukses, hapus job dan hapus proses update job
	// err = ctrl.Svc.DeleteJob(jobID)
	// if err != nil {
	// 	statusCode, resp := helper.GetError(err.Error())
	// 	return c.Status(statusCode).JSON(resp)
	// }

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
	jobs, err := ctrl.Svc.GetJobs()
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"success",
		jobs,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
