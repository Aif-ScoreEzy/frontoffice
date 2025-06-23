package member

import (
	"encoding/json"
	"errors"
	"fmt"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/internal/apperror"
	"front-office/pkg/core/role"
	"io"
	"net/http"
	"time"
)

func NewService(repo Repository, roleRepo role.Repository) Service {
	return &service{
		repo,
		roleRepo,
	}
}

type service struct {
	repo     Repository
	roleRepo role.Repository
}

type Service interface {
	GetMemberBy(query *FindUserQuery) (*MstMember, error)
	GetMemberList(filter *MemberFilter) ([]*MstMember, *model.Meta, error)
	UpdateProfile(id, oldEmail string, req *UpdateProfileRequest) error
	UploadProfileImage(id string, filename *string) error
	UpdateMemberById(id string, req *UpdateUserRequest) error
	DeleteMemberById(id string) error
}

func (s *service) GetMemberBy(query *FindUserQuery) (*MstMember, error) {
	member, err := s.repo.CallGetMemberAPI(query)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to get member")
	}

	return member, nil
}

func (svc *service) GetMemberList(filter *MemberFilter) ([]*MstMember, *model.Meta, error) {
	if filter.RoleName != "" {
		roles, err := svc.roleRepo.CallGetRolesAPI(role.RoleFilter{
			Name: filter.RoleName,
		})
		if err != nil {
			return nil, nil, apperror.MapRepoError(err, "failed to fetch role")
		}

		if len(roles) == 0 {
			return nil, nil, apperror.NotFound("role not found")
		}

		filter.RoleID = fmt.Sprintf("%v", roles[0].RoleId)
	}

	users, meta, err := svc.repo.CallGetMemberListAPI(filter)
	if err != nil {
		return nil, nil, err
	}

	return users, meta, nil
}

func (s *service) UpdateProfile(id, oldEmail string, req *UpdateProfileRequest) error {
	updateUser := map[string]interface{}{}

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		updateUser["email"] = *req.Email
	}

	updateUser["updated_at"] = time.Now()

	return s.repo.CallUpdateMemberAPI(id, updateUser)
}

func (s *service) UploadProfileImage(id string, filename *string) error {
	updateUser := map[string]interface{}{}

	if filename != nil {
		updateUser["image"] = *filename
	}

	updateUser["updated_at"] = time.Now()

	return s.repo.CallUpdateMemberAPI(id, updateUser)
}

func (s *service) UpdateMemberById(id string, req *UpdateUserRequest) error {
	updateUser := map[string]interface{}{}
	currentTime := time.Now()

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		member, err := s.GetMemberBy(&FindUserQuery{
			Email: *req.Email,
		})
		if err != nil {
			return apperror.MapRepoError(err, "failed to get member")
		}

		if member.MemberId != 0 {
			return errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	if req.RoleId != nil {
		_, err := s.roleRepo.CallGetRoleAPI(*req.RoleId)
		if err != nil {
			return apperror.MapRepoError(err, "failed to fetch role")
		}

		updateUser["role_id"] = *req.RoleId
	}

	if req.Active != nil {
		updateUser["active"] = *req.Active
	}

	updateUser["updated_at"] = currentTime

	return s.repo.CallUpdateMemberAPI(id, updateUser)
}

func (s *service) DeleteMemberById(id string) error {
	err := s.repo.CallDeleteMemberAPI(id)
	if err != nil {
		return apperror.MapRepoError(err, "failed to delete member")
	}

	return nil
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
