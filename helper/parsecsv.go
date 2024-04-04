package helper

import (
	"encoding/csv"
	"mime/multipart"
)

func ParseCSVFile(file *multipart.FileHeader) ([][]string, int, error) {
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

	return csvData, len(csvData), nil
}
