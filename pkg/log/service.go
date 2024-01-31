package log

import (
	"encoding/json"
	"front-office/common/constant"
	"front-office/common/model"
	"io"
	"net/http"
	"os"
)

func GetAllLogTransSvc() (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := os.Getenv("AIFCORE_HOST") + os.Getenv("GET_ALL_LOG")

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResp)

	return dataResp, response.StatusCode, nil
}
