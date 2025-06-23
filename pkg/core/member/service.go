package member

import (
	"errors"
	"fmt"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
	"time"

	"github.com/rs/zerolog/log"
)

func NewService(repo Repository, roleRepo role.Repository, operationRepo operation.Repository) Service {
	return &service{
		repo,
		roleRepo,
		operationRepo,
	}
}

type service struct {
	repo          Repository
	roleRepo      role.Repository
	operationRepo operation.Repository
}

type Service interface {
	GetMemberBy(query *FindUserQuery) (*MstMember, error)
	GetMemberList(filter *MemberFilter) ([]*MstMember, *model.Meta, error)
	UpdateProfile(userrId string, roleId uint, req *UpdateProfileRequest) (*userUpdateResponse, error)
	UploadProfileImage(id string, filename *string) (*userUpdateResponse, error)
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

func (svc *service) UpdateProfile(userId string, roleId uint, req *UpdateProfileRequest) (*userUpdateResponse, error) {
	user, err := svc.repo.CallGetMemberAPI(&FindUserQuery{Id: userId})
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch member")
	}
	if user.MemberId == 0 {
		return nil, apperror.NotFound(constant.UserNotFound)
	}

	updateFields := make(map[string]interface{})
	shouldSendEmailConfirmation := false
	var newEmail string

	if req.Name != nil {
		updateFields["name"] = *req.Name
	}

	if req.Email != nil {
		if roleId == uint(memberRoleID) {
			return nil, apperror.Unauthorized("you are not allowed to update email")
		}

		existing, err := svc.repo.CallGetMemberAPI(&FindUserQuery{Email: *req.Email})
		if err != nil {
			return nil, apperror.MapRepoError(err, "failed to check existing email")
		}
		if existing.MemberId != 0 {
			return nil, apperror.Conflict(constant.EmailAlreadyExists)
		}

		updateFields["email"] = *req.Email
		shouldSendEmailConfirmation = true
		newEmail = *req.Email
	}

	updateFields["updated_at"] = time.Now()

	if err := svc.repo.CallUpdateMemberAPI(userId, updateFields); err != nil {
		return nil, apperror.MapRepoError(err, "failed to update member")
	}

	if shouldSendEmailConfirmation {
		if err := mailjet.SendConfirmationEmailUserEmailChangeSuccess(user.Name, user.Email, newEmail, helper.FormatWIB(time.Now())); err != nil {
			return nil, apperror.Internal("failed to send email confirmation", err)
		}
		user.Email = newEmail
	}

	if err := svc.operationRepo.AddLogOperation(&operation.AddLogRequest{
		MemberId:  user.MemberId,
		CompanyId: user.CompanyId,
		Action:    constant.EventUpdateProfile,
	}); err != nil {
		log.Warn().Err(err).Msg("failed to log profile update event")
	}

	return &userUpdateResponse{
		Id:        user.MemberId,
		Name:      user.Name,
		Email:     user.Email,
		Active:    user.Active,
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
	}, nil
}

func (svc *service) UploadProfileImage(userId string, filename *string) (*userUpdateResponse, error) {
	user, err := svc.repo.CallGetMemberAPI(&FindUserQuery{Id: userId})
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch member")
	}
	if user.MemberId == 0 {
		return nil, apperror.NotFound(constant.UserNotFound)
	}

	updateFields := make(map[string]interface{})

	if filename != nil {
		updateFields["image"] = *filename
	}

	updateFields["updated_at"] = time.Now()

	if err := svc.repo.CallUpdateMemberAPI(userId, updateFields); err != nil {
		return nil, apperror.MapRepoError(err, "failed to update member")
	}

	if err := svc.operationRepo.AddLogOperation(&operation.AddLogRequest{
		MemberId:  user.MemberId,
		CompanyId: user.CompanyId,
		Action:    constant.EventUpdateProfile,
	}); err != nil {
		log.Warn().Err(err).Msg("failed to log upload profile photo event")
	}

	return &userUpdateResponse{
		Id:        user.MemberId,
		Name:      user.Name,
		Email:     user.Email,
		Active:    user.Active,
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
	}, nil
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
