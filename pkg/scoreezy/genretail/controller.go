package genretail

import (
	"encoding/csv"
	"fmt"
	"front-office/helper"
	"front-office/pkg/core/grade"
	"front-office/pkg/core/log/operation"
	"io"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"front-office/common/constant"
)

func NewController(
	service Service,
	svcGrading grade.Service,
	svcLogOperation operation.Service,
) Controller {
	return &controller{
		Svc:             service,
		SvcGrading:      svcGrading,
		SvcLogOperation: svcLogOperation,
	}
}

type controller struct {
	Svc             Service
	SvcGrading      grade.Service
	SvcLogOperation operation.Service
}

type Controller interface {
	RequestScore(c *fiber.Ctx) error
	DownloadCSV(c *fiber.Ctx) error
	UploadCSV(c *fiber.Ctx) error
	GetBulkSearch(c *fiber.Ctx) error
}

func (ctrl *controller) RequestScore(c *fiber.Ctx) error {
	req := c.Locals(constant.Request).(*GenRetailRequest)
	apiKey := c.Get(constant.XAPIKey)
	companyId := c.Locals(constant.CompanyId)

	currentUserId, err := helper.InterfaceToUint(c.Locals(constant.UserId))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	companyIdUint, err := helper.InterfaceToUint(c.Locals(constant.CompanyId))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	// make sure parameter settings are set
	productSlug := constant.SlugGenRetailV3
	result, err := ctrl.SvcGrading.GetGrades(productSlug, fmt.Sprintf("%v", companyId))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	if len(result.Grades) < 1 {
		statusCode, resp := helper.GetError(constant.ParamSettingIsNotSet)
		return c.Status(statusCode).JSON(resp)
	}

	genRetailResponse, errRequest := ctrl.Svc.GenRetailV3(req, apiKey)
	if errRequest != nil {
		statusCode, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if genRetailResponse.StatusCode >= 400 {
		dataReturn := GenRetailV3ClientReturnError{
			Message:      genRetailResponse.Message,
			ErrorMessage: genRetailResponse.ErrorMessage,
			Data:         genRetailResponse.Data,
		}

		return c.Status(genRetailResponse.StatusCode).JSON(dataReturn)
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  currentUserId,
		CompanyId: companyIdUint,
		Action:    constant.EventCalculateScore,
	}

	err = ctrl.SvcLogOperation.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for calculate score")
	}

	resp := GenRetailV3ClientReturnSuccess{
		Message: genRetailResponse.Message,
		Success: true,
		Data:    genRetailResponse.Data,
	}

	return c.Status(genRetailResponse.StatusCode).JSON(resp)
}

func (ctrl *controller) DownloadCSV(c *fiber.Ctx) error {
	opsi := c.Params("opsi")

	filePath := fmt.Sprintf("./public/bulk_template/%s.csv", opsi)

	return c.SendFile(filePath)
}

func (ctrl *controller) UploadCSV(c *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", c.Locals(constant.UserId))
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	tempType := fmt.Sprintf("%v", c.Locals("tempType"))
	apiKey := c.Get(constant.XAPIKey)

	// Get the file from the form data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		statusCode, resp := helper.GetError(constant.ErrorGettingFile)
		return c.Status(statusCode).JSON(resp)
	}

	file, err := fileHeader.Open()
	if err != nil {
		statusCode, resp := helper.GetError(constant.ErrorOpeningFile)
		return c.Status(statusCode).JSON(resp)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	header, err := reader.Read()
	if err != nil {
		statusCode, resp := helper.GetError(constant.ErrorReadingCSV)
		return c.Status(statusCode).JSON(resp)
	}

	// Process the header (first line)
	var validHeaderTemplate []string
	if tempType == "personal" {
		validHeaderTemplate = append(validHeaderTemplate, "loan_no", "name", "nik", "phone_number")
	} else {
		validHeaderTemplate = append(validHeaderTemplate, "company_id", "company_name", "npwp_company", "phone_number")
	}

	for _, v := range header {
		isValidHeader := helper.IsValidTemplateHeader(validHeaderTemplate, v)

		if !isValidHeader {
			statusCode, resp := helper.GetError(constant.HeaderTemplateNotValid)
			return c.Status(statusCode).JSON(resp)
		}
	}

	storeData := []BulkSearchRequest{}
	// Iterate over CSV records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			statusCode, resp := helper.GetError(constant.ErrorReadingCSVRecords)
			return c.Status(statusCode).JSON(resp)
		}

		// Process the CSV record
		insertNew := BulkSearchRequest{}
		for _, v := range record {
			fmt.Println("v: ", v)
			insertNew.LoanNo = record[0]
			insertNew.Name = record[1]
			insertNew.NIK = record[2]
			insertNew.PhoneNumber = record[3]
		}
		storeData = append(storeData, insertNew)
	}

	processInsert := ctrl.Svc.BulkSearchUploadSvc(storeData, tempType, apiKey, userId, companyId)

	if processInsert != nil {
		statusCode, resp := helper.GetError(constant.ErrorUploadDataCSV)
		return c.Status(statusCode).JSON(resp)
	}

	bulkSearch, err := ctrl.Svc.GetBulkSearchSvc(uint(tierLevel), userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetTotalDataBulk(uint(tierLevel), userId, companyId)

	fullResponsePage := map[string]interface{}{
		"total_data": totalData,
		"data":       bulkSearch,
	}

	resp := helper.ResponseSuccess(
		"succeed to upload data",
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetBulkSearch(c *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", c.Locals(constant.UserId))
	companyId := fmt.Sprintf("%v", c.Locals(constant.CompanyId))
	tierLevel, _ := strconv.ParseUint(fmt.Sprintf("%v", c.Locals("tierLevel")), 10, 64)
	// find user loggin detail

	bulkSearch, err := ctrl.Svc.GetBulkSearchSvc(uint(tierLevel), userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetTotalDataBulk(uint(tierLevel), userId, companyId)

	fullResponsePage := map[string]interface{}{
		"total_data": totalData,
		"data":       bulkSearch,
	}

	resp := helper.ResponseSuccess(
		"succeed to get bulk search data",
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
