package helper

import (
	"encoding/csv"
	"errors"
	"front-office/common/constant"
	"log"
	"mime/multipart"
)

func ParseCSVFile(file *multipart.FileHeader, expectedHeaders []string) ([][]string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := fileContent.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	reader := csv.NewReader(fileContent)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := csvData[0]
	for i, expectedHeader := range expectedHeaders {
		if header[i] != expectedHeader {
			return nil, errors.New(constant.HeaderTemplateNotValid)
		}
	}

	// totalData := len(csvData) - 1

	return csvData, nil
}
