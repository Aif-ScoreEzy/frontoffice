package log

import (
	"encoding/json"
	"front-office/common/constant"
	"front-office/common/model"
	"io"
	"net/http"
	"os"
)

func GetAllLogTransByDateSvc(companyID, date string) (*model.AifResponse, int, error) {
	var dataResp *model.AifResponse
	url := os.Getenv("AIFCORE_HOST") + os.Getenv("GET_ALL_LOG_BY_DATE")

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := request.URL.Query()
	q.Add("company_id", companyID)
	q.Add("date", date)
	request.URL.RawQuery = q.Encode()

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
