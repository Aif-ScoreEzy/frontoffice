package phonelivestatus

import (
	"encoding/json"
	"errors"
	"front-office/app/config"
	"io"
	"net/http"
)

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{
		Cfg:  cfg,
		Repo: repo,
	}
}

type service struct {
	Cfg  *config.Config
	Repo Repository
}

type Service interface {
	GetPhoneLiveStatusJobAPI(filter *PhoneLiveStatusFilter) (*APIResponse[JobListResponse], error)
	ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error
}

func (svc *service) GetPhoneLiveStatusJobAPI(filter *PhoneLiveStatusFilter) (*APIResponse[JobListResponse], error) {
	response, err := svc.Repo.CallGetPhoneLiveStatusJobAPI(filter)
	if err != nil {
		return nil, err
	}

	result, err := parseGenericResponse[JobListResponse](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error {
	_, err := svc.Repo.CallPhoneLiveStatusAPI(memberId, companyId, req)
	if err != nil {
		return err
	}

	return nil
}

func parseGenericResponse[T any](response *http.Response) (*APIResponse[T], error) {
	var apiResponse APIResponse[T]

	if response == nil {
		return nil, errors.New("nil response")
	}

	dataBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, &apiResponse); err != nil {
		return nil, err
	}

	apiResponse.StatusCode = response.StatusCode
	return &apiResponse, nil
}
