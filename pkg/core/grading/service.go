package grading

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
	GetGradings(companyId string) (*AifResponse, error)
}

func (svc *service) GetGradings(companyId string) (*AifResponse, error) {
	response, err := svc.Repo.GetGradeList(companyId)
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
