package user

import (
	"encoding/json"
	"errors"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
	"io"
	"strconv"
	"time"
)

func NewService(repo Repository, repoRole role.Repository) Service {
	return &service{Repo: repo, RepoRole: repoRole}
}

type service struct {
	Repo     Repository
	RepoRole role.Repository
}

type Service interface {
	FindUserByEmailSvc(email string) (*User, error)
	FindUserByKeySvc(key string) (*User, error)
	FindUserByIdSvc(id string) (*User, error)
	FindUserByIdAndCompanyIdSvc(id, companyId string) (*User, error)
	UpdateProfileSvc(req *UpdateProfileRequest, user *User) (*User, error)
	UploadProfileImageSvc(user *User, filename *string) (*User, error)
	UpdateUserByIdSvc(req *UpdateUserRequest, user *User) (*User, error)
	GetAllUsersSvc(limit, page, keyword, roleId, status, startDate, endDate, companyId string) ([]GetUsersResponse, error)
	GetTotalDataSvc(keyword, roleId, active, startDate, endDate, companyId string) (int64, error)
	DeleteUserByIdSvc(id string) error
	FindUserAifCore(query *FindUserQuery) (*FindUserAifCoreResponse, error)
	UpdateUserByIdAifCore(req *UpdateUserRequest, memberId uint) error
}

func (svc *service) FindUserByEmailSvc(email string) (*User, error) {
	user, err := svc.Repo.FindOneByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) FindUserByKeySvc(key string) (*User, error) {
	user, err := svc.Repo.FindOneByKey(key)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) FindUserByIdSvc(id string) (*User, error) {
	user, err := svc.Repo.FindOneByUserId(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (svc *service) FindUserByIdAndCompanyIdSvc(id, companyId string) (*User, error) {
	user, err := svc.Repo.FindOneByUserIdAndCompanyId(id, companyId)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (svc *service) UpdateProfileSvc(req *UpdateProfileRequest, user *User) (*User, error) {
	updateUser := map[string]interface{}{}

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		result, _ := svc.Repo.FindOneByUserIdAndCompanyId(user.Id, user.CompanyId)
		if result.Role.TierLevel == 2 {
			return nil, errors.New(constant.RequestProhibited)
		}

		result, _ = svc.Repo.FindOneByEmail(*req.Email)
		if result != nil {
			return nil, errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	updateUser["updated_at"] = time.Now()

	user, err := svc.Repo.UpdateOneById(updateUser, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) UploadProfileImageSvc(user *User, filename *string) (*User, error) {
	updateUser := map[string]interface{}{}

	if filename != nil {
		updateUser["image"] = *filename
	}

	updateUser["updated_at"] = time.Now()

	user, err := svc.Repo.UpdateOneById(updateUser, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) UpdateUserByIdSvc(req *UpdateUserRequest, user *User) (*User, error) {
	updateUser := map[string]interface{}{}
	oldEmail := user.Email
	currentTime := time.Now()

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		userExists, _ := svc.Repo.FindOneByEmail(*req.Email)
		if userExists != nil {
			return nil, errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	if req.RoleId != nil {
		role, err := svc.RepoRole.FindOneById(*req.RoleId)
		if role == nil {
			return nil, errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, err
		}

		updateUser["role_id"] = *req.RoleId
	}

	if req.Active != nil {
		if *req.Active {
			updateUser["status"] = "active"
			updateUser["active"] = true
		} else {
			updateUser["status"] = "inactive"
			updateUser["active"] = false
		}
	}

	if req.Status != nil {
		updateUser["status"] = *req.Status
	}

	updateUser["updated_at"] = currentTime

	updatedUser, err := svc.Repo.UpdateOneById(updateUser, user)
	if err != nil {
		return nil, err
	}

	formattedTime := helper.FormatWIB(currentTime)

	if oldEmail != updatedUser.Email {
		err := mailjet.SendConfirmationEmailUserEmailChangeSuccess(updatedUser.Name, oldEmail, *req.Email, formattedTime)
		if err != nil {
			return nil, err
		}
	}

	return updatedUser, nil
}

func (svc *service) GetAllUsersSvc(limit, page, keyword, roleId, status, startDate, endDate, companyId string) ([]GetUsersResponse, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	if status != "" && (status != "active" && status != "inactive" && status != "resend" && status != "pending") {
		return nil, errors.New(constant.InvalidStatusValue)
	}

	var startTime, endTime string
	layoutPostgreSQLDate := "2006-01-02"
	if startDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, startDate)
		if err != nil {
			return nil, errors.New(constant.InvalidDateFormat)
		}

		startTime = helper.FormatStartTimeForSQL(startDate)

		if endDate == "" {
			endTime = helper.FormatEndTimeForSQL(startDate)
		}
	}

	if endDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, endDate)
		if err != nil {
			return nil, errors.New(constant.InvalidDateFormat)
		}

		endTime = helper.FormatEndTimeForSQL(endDate)
	}

	users, err := svc.Repo.FindAll(intLimit, offset, keyword, roleId, status, startTime, endTime, companyId)
	if err != nil {
		return nil, err
	}

	var responseUsers []GetUsersResponse
	for _, user := range users {
		responseUser := GetUsersResponse{
			Id:         user.Id,
			Name:       user.Name,
			Email:      user.Email,
			Status:     user.Status,
			Active:     user.Active,
			IsVerified: user.IsVerified,
			CompanyId:  user.CompanyId,
			Role:       user.Role,
			CreatedAt:  user.CreatedAt,
		}
		responseUsers = append(responseUsers, responseUser)
	}

	return responseUsers, nil
}

func (svc *service) GetTotalDataSvc(keyword, roleId, active, startDate, endDate, companyId string) (int64, error) {
	var startTime, endTime string
	layoutPostgreSQLDate := "2006-01-02"
	if startDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, startDate)
		if err != nil {
			return 0, errors.New(constant.InvalidDateFormat)
		}

		startTime = helper.FormatStartTimeForSQL(startDate)

		if endDate == "" {
			endTime = helper.FormatEndTimeForSQL(startDate)
		}
	}

	if endDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, endDate)
		if err != nil {
			return 0, errors.New(constant.InvalidDateFormat)
		}

		endTime = helper.FormatEndTimeForSQL(endDate)
	}

	count, err := svc.Repo.GetTotalData(keyword, roleId, active, startTime, endTime, companyId)
	return count, err
}

func (svc *service) DeleteUserByIdSvc(id string) error {
	err := svc.Repo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) FindUserAifCore(query *FindUserQuery) (*FindUserAifCoreResponse, error) {
	res, err := svc.Repo.FindOneAifCore(query)
	if err != nil {
		return nil, err
	}

	var baseResponseSuccess *FindUserAifCoreResponse
	if res != nil {
		dataBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		json.Unmarshal(dataBytes, &baseResponseSuccess)
		baseResponseSuccess.StatusCode = res.StatusCode
	}

	return baseResponseSuccess, nil
}

func (svc *service) UpdateUserByIdAifCore(req *UpdateUserRequest, memberId uint) error {
	updateUser := map[string]interface{}{}

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		userExists, _ := svc.Repo.FindOneByEmail(*req.Email)
		if userExists != nil {
			return errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	if req.RoleId != nil {
		role, err := svc.RepoRole.FindOneById(*req.RoleId)
		if role == nil {
			return errors.New(constant.DataNotFound)
		} else if err != nil {
			return err
		}

		updateUser["role_id"] = *req.RoleId
	}

	if req.Active != nil {
		if *req.Active {
			updateUser["status"] = "active"
			updateUser["active"] = true
		} else {
			updateUser["status"] = "inactive"
			updateUser["active"] = false
		}
	}

	if req.Status != nil {
		updateUser["status"] = *req.Status
	}

	_, err := svc.Repo.UpdateOneByIdAifCore(updateUser, memberId)
	if err != nil {
		return err
	}

	return nil

	// currentTime := time.Now()
	// formattedTime := helper.FormatWIB(currentTime)

	// if oldEmail != updatedUser.Email {
	// 	err := mailjet.SendConfirmationEmailUserEmailChangeSuccess(updatedUser.Name, oldEmail, *req.Email, formattedTime)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
}
