package transaction

import (
	"encoding/json"
	"io"
	"net/http"
)

func (svc *service) GetLogScoreezy() (*AifResponse, int, error) {
	response, err := svc.repo.CallLogScoreezyAPI()
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetLogScoreezyByDate(companyId, date string) (*AifResponse, int, error) {
	response, err := svc.repo.CallLogScoreezyByDateAPI(companyId, date)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetLogScoreezyByRangeDate(startDate, endDate, companyId, page string) (*AifResponse, int, error) {
	response, err := svc.repo.CallLogScoreezyByRangeDateAPI(companyId, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func (svc *service) GetLogScoreezyByMonth(companyId, month string) (*AifResponse, int, error) {
	response, err := svc.repo.CallLogScoreezyByMonthAPI(companyId, month)
	if err != nil {
		return nil, 0, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, 0, err
	}

	return result, response.StatusCode, nil
}

func parseResponse(response *http.Response) (*AifResponse, error) {
	var baseResponse *AifResponse

	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponse); err != nil {
			return nil, err
		}
	}

	return baseResponse, nil
}
