package genretail

import (
	"encoding/csv"
	"fmt"
	"front-office/helper"
	"io"

	"front-office/pkg/grading"
	"front-office/pkg/user"

	"github.com/gofiber/fiber/v2"

	"front-office/common/constant"
)

func RequestScore(c *fiber.Ctx) error {
	req := c.Locals("request").(*GenRetailRequest)
	apiKey := c.Get("X-API-KEY")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	// make sure parameter settings are set
	gradings, _ := grading.GetGradingsSvc(companyID)
	if len(gradings) < 1 {
		statusCode, resp := helper.GetError(constant.ParamSettingIsNotSet)
		return c.Status(statusCode).JSON(resp)
	}

	genRetailResponse, errRequest := GenRetailV3(req, apiKey)
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

	resp := GenRetailV3ClientReturnSuccess{
		Message: genRetailResponse.Message,
		Success: true,
		Data:    genRetailResponse.Data,
	}

	return c.Status(genRetailResponse.StatusCode).JSON(resp)
}

func DownloadCSV(c *fiber.Ctx) error {
	opsi := c.Params("opsi")

	filePath := fmt.Sprintf("./public/bulk_template/%s.csv", opsi)

	return c.SendFile(filePath)
}

func UploadCSV(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	tempType := fmt.Sprintf("%v", c.Locals("tempType"))
	apiKey := c.Get("X-API-KEY")
	userDetails := c.Locals("userDetails").(*user.User)

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
	fmt.Println("CSV Header:", header)

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

	processInsert := BulkSearchUploadSvc(storeData, tempType, apiKey, userID, companyID)

	if processInsert != nil {
		statusCode, resp := helper.GetError(constant.ErrorUploadDataCSV)
		return c.Status(statusCode).JSON(resp)
	}

	bulkSearch, err := GetBulkSearchSvc(userDetails.Role.TierLevel, userDetails.ID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := GetTotalDataBulk(userDetails.Role.TierLevel, userDetails.ID, companyID)

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

func GetBulkSearch(c *fiber.Ctx) error {
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	userDetails := c.Locals("userDetails").(*user.User)
	// find user loggin detail

	fmt.Println("userDetails.Role.TierLevel ", userDetails.Role.TierLevel)
	fmt.Println("userID ", userDetails.ID)

	bulkSearch, err := GetBulkSearchSvc(userDetails.Role.TierLevel, userDetails.ID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := GetTotalDataBulk(userDetails.Role.TierLevel, userDetails.ID, companyID)

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
