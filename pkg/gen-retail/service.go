package genretail

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func GenRetailV3(requestData *GenRetailRequest, apiKey string) (*GenRetailV3ModelResponse, error) {
	var dataResponse *GenRetailV3ModelResponse

	URL := os.Getenv("GEN_RETAIL_HOST") + os.Getenv("GEN_RETAIL_V3")

	requestByte, _ := json.Marshal(requestData)
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestByte))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", apiKey)

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
