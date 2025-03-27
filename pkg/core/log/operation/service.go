package operation

import (
	"encoding/json"
	"io"
	"net/http"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	GetLogOperations(filter *LogOperationFilter) (*AifResponse, error)
	GetByRange(filter *LogRangeFilter) (*AifResponse, error)
	AddLogOperation(req *AddLogRequest) (*AifResponse, error)
}

func (svc *service) GetLogOperations(filter *LogOperationFilter) (*AifResponse, error) {
	response, err := svc.Repo.FetchLogOperations(filter)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) GetByRange(filter *LogRangeFilter) (*AifResponse, error) {
	response, err := svc.Repo.FetchByRange(filter)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) AddLogOperation(req *AddLogRequest) (*AifResponse, error) {
	response, err := svc.Repo.AddLogOperation(req)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
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
