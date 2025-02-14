package member

import (
	"encoding/json"
	"errors"
	"fmt"
	"front-office/common/constant"
	"front-office/pkg/core/role"
	"io"
	"net/http"
	"time"
)

func NewService(repo Repository, roleSvc role.Service) Service {
	return &service{
		Repo:    repo,
		RoleSvc: roleSvc,
	}
}

type service struct {
	Repo    Repository
	RoleSvc role.Service
}

type Service interface {
	GetMemberBy(query *FindUserQuery) (*AifResponse, error)
	GetMemberList(filter *MemberFilter) (*AifResponseWithMultipleData, error)
	UpdateProfile(id, oldEmail string, req *UpdateProfileRequest) (*AifResponse, error)
	UploadProfileImage(id string, filename *string) (*AifResponse, error)
	UpdateMemberByIdSvc(id string, req *UpdateUserRequest) (*AifResponse, error)
	DeleteMemberById(id string) (*AifResponse, error)
}

func (s *service) GetMemberBy(query *FindUserQuery) (*AifResponse, error) {
	response, err := s.Repo.GetMemberBy(query)
	if err != nil {
		return nil, err
	}

	return s.parseSingleResponse(response)
}

func (s *service) GetMemberList(filter *MemberFilter) (*AifResponseWithMultipleData, error) {
	response, err := s.Repo.GetMemberList(filter)
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

func (s *service) UpdateMemberByIdSvc(id string, req *UpdateUserRequest) (*AifResponse, error) {
	updateUser := map[string]interface{}{}
	currentTime := time.Now()

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		result, _ := s.GetMemberBy(&FindUserQuery{
			Email: *req.Email,
		})

		if result.Data.MemberId != 0 {
			return nil, errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	if req.RoleId != nil {
		role, err := s.RoleSvc.GetRoleById(*req.RoleId)
		if err != nil {
			return nil, err
		}

		if role.Data.RoleId == 0 {
			return nil, errors.New(constant.DataNotFound)
		}

		updateUser["role_id"] = *req.RoleId
	}

	if req.Active != nil {
		updateUser["active"] = *req.Active
	}

	updateUser["updated_at"] = currentTime

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
