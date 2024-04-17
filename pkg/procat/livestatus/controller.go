package livestatus

import (
	"front-office/helper"

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
}

func (ctrl *controller) BulkSearch(c *fiber.Ctx) error {
	apiKey := c.Get("X-AIF-KEY")

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

	// var liveStatusResponse *LiveStatusResponse
	// var liveStatusResponses []*LiveStatusResponse
	// if countJobDetail != 0 {
	jobDetails, err := ctrl.Svc.GetJobDetails(jobID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	batchSize := 10
	// jobIDStr := strconv.FormatUint(uint64(jobID), 10)
	liveStatusResponses, err := ctrl.Svc.ProcessJobDetails(apiKey, jobID, jobDetails, batchSize)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	// for _, jobDetail := range jobDetails {
	// 	request := &LiveStatusRequest{
	// 		PhoneNumber: jobDetail.PhoneNumber,
	// 		TrxID:       jobIDStr,
	// 	}

	// 	liveStatusResponse, err = ctrl.Svc.CreateLiveStatus(request, apiKey)
	// 	if err != nil {
	// 		statusCode, resp := helper.GetError(err.Error())
	// 		return c.Status(statusCode).JSON(resp)
	// 	}

	// 	liveStatusResponses = append(liveStatusResponses, liveStatusResponse)

	// 	err = ctrl.Svc.DeleteJobDetail(jobDetail.ID)
	// 	if err != nil {
	// 		statusCode, resp := helper.GetError(err.Error())
	// 		return c.Status(statusCode).JSON(resp)
	// 	}
	// }
	// }

	// if liveStatusResponse.StatusCode >= 400 {
	// 	dataReturn := helper.BaseResponseFailed{
	// 		Message: liveStatusResponse.Message,
	// 	}

	// 	return c.Status(liveStatusResponse.StatusCode).JSON(dataReturn)
	// }

	resp := helper.ResponseSuccess(
		"success",
		liveStatusResponses,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
