package genretail

import (
	"bytes"
	"encoding/json"
	"front-office/common/constant"
	"io"
	"net/http"
	"os"
)

func GenRetailV3(requestData *GenRetailRequest, apiKey string) (*GenRetailV3ModelResponse, error) {
	var dataResponse *GenRetailV3ModelResponse

	url := os.Getenv("AIFCORE_HOST") + os.Getenv("GEN_RETAIL_V3")

	requestByte, _ := json.Marshal(requestData)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestByte))

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	request.Header.Set(constant.XAPIKey, apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResponse)
	dataResponse.StatusCode = response.StatusCode

	return dataResponse, nil
}
