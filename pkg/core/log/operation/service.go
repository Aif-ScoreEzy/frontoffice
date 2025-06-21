package operation

import (
	"encoding/json"
	"front-office/internal/apperror/mapper"
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
	GetLogOperations(filter *LogOperationFilter) (log *LogOperation, err error)
	GetByRange(filter *LogRangeFilter) (*AifResponse, error)
	AddLogOperation(req *AddLogRequest) error
}

func (svc *service) GetLogOperations(filter *LogOperationFilter) (*LogOperation, error) {
	log, err := svc.Repo.FetchLogOperations(filter)
	if err != nil {
		return nil, mapper.MapRepoError(err, "failed to fetch log operations")
	}

	return log, nil
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

func (svc *service) AddLogOperation(req *AddLogRequest) error {
	return svc.Repo.AddLogOperation(req)
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
