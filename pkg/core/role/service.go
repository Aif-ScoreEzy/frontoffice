package role

import (
	"encoding/json"
	"fmt"
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
	GetRoleById(id string) (*AifResponse, error)
	GetAllRoles() (*AifResponseWithMultipleData, error)
}

func (s *service) GetRoleById(id string) (*AifResponse, error) {
	res, err := s.Repo.FindOneById(id)
	if err != nil {
		return nil, err
	}

	return parseSingleResponse(res)
}

func (s *service) GetAllRoles() (*AifResponseWithMultipleData, error) {
	res, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	return parseMultipleResponse(res)
}

func parseResponse(response *http.Response, result interface{}) error {
	if response == nil {
		return fmt.Errorf("response is nil")
	}
	defer response.Body.Close()

	dataByte, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(dataByte, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func parseSingleResponse(response *http.Response) (*AifResponse, error) {
	var baseResponse AifResponse
	if err := parseResponse(response, &baseResponse); err != nil {
		return nil, err
	}

	return &baseResponse, nil
}

func parseMultipleResponse(response *http.Response) (*AifResponseWithMultipleData, error) {
	var baseResponse AifResponseWithMultipleData
	if err := parseResponse(response, &baseResponse); err != nil {
		return nil, err
	}

	return &baseResponse, nil
}
