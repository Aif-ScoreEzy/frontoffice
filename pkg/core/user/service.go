package user

import (
	"errors"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
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
	FindUserByIDSvc(id string) (*User, error)
	FindUserByIDAndCompanyIDSvc(id, companyID string) (*User, error)
	UpdateProfileSvc(req *UpdateProfileRequest, user *User) (*User, error)
	UploadProfileImageSvc(user *User, filename *string) (*User, error)
	UpdateUserByIDSvc(req *UpdateUserRequest, user *User) (*User, error)
	GetAllUsersSvc(limit, page, keyword, roleID, status, startDate, endDate, companyID string) ([]GetUsersResponse, error)
	GetTotalDataSvc(keyword, roleID, active, startDate, endDate, companyID string) (int64, error)
	DeleteUserByIDSvc(id string) error
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

func (svc *service) FindUserByIDSvc(id string) (*User, error) {
	user, err := svc.Repo.FindOneByUserID(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (svc *service) FindUserByIDAndCompanyIDSvc(id, companyID string) (*User, error) {
	user, err := svc.Repo.FindOneByUserIDAndCompanyID(id, companyID)
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
		result, _ := svc.Repo.FindOneByUserIDAndCompanyID(user.ID, user.CompanyID)
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

	user, err := svc.Repo.UpdateOneByID(updateUser, user)
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

	user, err := svc.Repo.UpdateOneByID(updateUser, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) UpdateUserByIDSvc(req *UpdateUserRequest, user *User) (*User, error) {
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

	if req.RoleID != nil {
		role, err := svc.RepoRole.FindOneByID(*req.RoleID)
		if role == nil {
			return nil, errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, err
		}

		updateUser["role_id"] = *req.RoleID
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

	updatedUser, err := svc.Repo.UpdateOneByID(updateUser, user)
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

func (svc *service) GetAllUsersSvc(limit, page, keyword, roleID, status, startDate, endDate, companyID string) ([]GetUsersResponse, error) {
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

	users, err := svc.Repo.FindAll(intLimit, offset, keyword, roleID, status, startTime, endTime, companyID)
	if err != nil {
		return nil, err
	}

	var responseUsers []GetUsersResponse
	for _, user := range users {
		responseUser := GetUsersResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Status:     user.Status,
			Active:     user.Active,
			IsVerified: user.IsVerified,
			CompanyID:  user.CompanyID,
			Role:       user.Role,
			CreatedAt:  user.CreatedAt,
		}
		responseUsers = append(responseUsers, responseUser)
	}

	return responseUsers, nil
}

func (svc *service) GetTotalDataSvc(keyword, roleID, active, startDate, endDate, companyID string) (int64, error) {
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

	count, err := svc.Repo.GetTotalData(keyword, roleID, active, startTime, endTime, companyID)
	return count, err
}

func (svc *service) DeleteUserByIDSvc(id string) error {
	err := svc.Repo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
