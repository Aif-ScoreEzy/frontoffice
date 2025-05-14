package phonelivestatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"front-office/app/config"
	"io"
	"mime/multipart"
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
	GetPhoneLiveStatusJob(filter *PhoneLiveStatusFilter) (*APIResponse[JobListResponse], error)
	GetPhoneLiveStatusDetails(filter *PhoneLiveStatusFilter) (*APIResponse[JobDetailsResponse], error)
	ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error
	BulkProcessPhoneLiveStatus(memberId, companyId string, fileHeader *multipart.FileHeader) error
}

func (svc *service) GetPhoneLiveStatusJob(filter *PhoneLiveStatusFilter) (*APIResponse[JobListResponse], error) {
	response, err := svc.Repo.CallGetPhoneLiveStatusJobAPI(filter)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("API error: %s, body: %s", response.Status, string(body))
	}

	return parseGenericResponse[JobListResponse](response)
}

func (svc *service) GetPhoneLiveStatusDetails(filter *PhoneLiveStatusFilter) (*APIResponse[JobDetailsResponse], error) {
	response, err := svc.Repo.CallGetJobDetailsAPI(filter)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("API error: %s, body: %s", response.Status, string(body))
	}

	return parseGenericResponse[JobDetailsResponse](response)
}

func (svc *service) ProcessPhoneLiveStatus(memberId, companyId string, req *PhoneLiveStatusRequest) error {
	_, err := svc.Repo.CallPhoneLiveStatusAPI(memberId, companyId, req)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) BulkProcessPhoneLiveStatus(memberId, companyId string, fileHeader *multipart.FileHeader) error {
	_, err := svc.Repo.CallBulkPhoneLiveStatusAPI(memberId, companyId, fileHeader)
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
