package helper

import (
	"encoding/csv"
	"errors"
	"front-office/common/constant"
	"mime/multipart"
)

func ParseCSVFile(file *multipart.FileHeader, expectedHeaders []string) ([][]string, int, error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, 0, err
	}
	defer fileContent.Close()

	reader := csv.NewReader(fileContent)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, 0, err
	}

	header := csvData[0]
	for i, expectedHeader := range expectedHeaders {
		if header[i] != expectedHeader {
			return nil, 0, errors.New(constant.HeaderTemplateNotValid)
		}
	}

	return csvData, len(csvData), nil
}
