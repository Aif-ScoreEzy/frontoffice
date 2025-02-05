package member

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewService(repo Repository) Service {
	return &service{
		Repo: repo,
	}
}

type service struct {
	Repo Repository
}

type Service interface {
	GetMemberBy(query *FindUserQuery) (*AifResponse, error)
	GetMemberList(companyId string) (*AifResponseWithMultipleData, error)
	UpdateProfile(id, oldEmail string, req *UpdateProfileRequest) (*AifResponse, error)
	UploadProfileImage(id string, filename *string) (*AifResponse, error)
	DeleteMemberById(id string) (*AifResponse, error)
}

func (s *service) GetMemberBy(query *FindUserQuery) (*AifResponse, error) {
	response, err := s.Repo.GetMemberBy(query)
	if err != nil {
		return nil, err
	}

	return s.parseSingleResponse(response)
}

func (s *service) GetMemberList(companyId string) (*AifResponseWithMultipleData, error) {
	response, err := s.Repo.GetMemberList(companyId)
	if err != nil {
		return nil, err
	}

	return s.parseMultipleResponse(response)
}

func (s *service) UpdateProfile(id, oldEmail string, req *UpdateProfileRequest) (*AifResponse, error) {
	updateUser := map[string]interface{}{}

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		updateUser["email"] = *req.Email
	}

	updateUser["updated_at"] = time.Now()

	response, err := s.Repo.UpdateOneById(id, updateUser)
	if err != nil {
		return nil, err
	}

	return s.parseSingleResponse(response)
}

func (s *service) UploadProfileImage(id string, filename *string) (*AifResponse, error) {
	updateUser := map[string]interface{}{}

	if filename != nil {
		updateUser["image"] = *filename
	}

	updateUser["updated_at"] = time.Now()

	response, err := s.Repo.UpdateOneById(id, updateUser)
	if err != nil {
		return nil, err
	}

	return s.parseSingleResponse(response)
}

func (s *service) DeleteMemberById(id string) (*AifResponse, error) {
	response, err := s.Repo.DeleteMemberById(id)
	if err != nil {
		return nil, err
	}

	return s.parseSingleResponse(response)
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

func (s *service) parseSingleResponse(response *http.Response) (*AifResponse, error) {
	var baseResponse AifResponse
	if err := parseResponse(response, &baseResponse); err != nil {
		return nil, err
	}

	return &baseResponse, nil
}

func (s *service) parseMultipleResponse(response *http.Response) (*AifResponseWithMultipleData, error) {
	var baseResponse AifResponseWithMultipleData
	if err := parseResponse(response, &baseResponse); err != nil {
		return nil, err
	}

	return &baseResponse, nil
}
