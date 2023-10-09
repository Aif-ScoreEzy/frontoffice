package user

import (
	"errors"
	"front-office/constant"
	"front-office/helper"
	"front-office/utility/mailjet"
	"strconv"

	"github.com/google/uuid"
)

func RegisterMemberSvc(req *RegisterMemberRequest, loggedUser *User) (*User, error) {
	userID := uuid.NewString()
	dataUser := &User{
		ID:         userID,
		Name:       req.Name,
		Email:      req.Email,
		Key:        helper.GenerateAPIKey(),
		RoleID:     req.RoleID,
		Active:     req.Active,
		IsVerified: true,
		CompanyID:  loggedUser.CompanyID,
	}

	password := helper.GeneratePassword()
	dataUser.Password = SetPassword(password)

	user, err := CreateMember(dataUser)
	if err != nil {
		return user, err
	}

	variables := map[string]interface{}{
		"email":    req.Email,
		"password": password,
	}

	err = mailjet.CreateMailjet(req.Email, 5082139, variables)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByEmailSvc(email string) (*User, error) {
	user, err := FindOneByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByKeySvc(key string) (*User, error) {
	user, err := FindOneByKey(key)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByIDSvc(id string) (*User, error) {
	user, err := FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func ActivateUserByKeySvc(key string) (*User, error) {
	user, err := UpdateOneByKey(key)
	if err != nil {
		return user, err
	}

	return user, nil
}

func DeactivateUserByEmailSvc(email string) (*User, error) {
	user, err := DeactiveOneByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUserByIDSvc(req *UpdateUserRequest, id string) (*User, error) {

	dataReq := &User{}

	if req.Name != "" {
		dataReq.Name = req.Name
	}
	if req.Email != "" {
		dataReq.Email = req.Email
	}
	if req.Phone != "" {
		dataReq.Phone = req.Phone
	}
	if req.CompanyID != "" {
		dataReq.CompanyID = req.CompanyID
	}
	if req.RoleID != "" {
		dataReq.RoleID = req.RoleID
	}

	user, err := UpdateOneByID(dataReq, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetAllUsersSvc(limit, page, keyword, roleID, active, startDate, endDate, companyID string) ([]GetUsersResponse, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	if active != "" {
		_, err := strconv.ParseBool(active)
		if err != nil {
			return nil, errors.New(constant.InvalidActiveValue)
		}
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

	users, err := FindAll(intLimit, offset, keyword, roleID, active, startTime, endTime, companyID)
	if err != nil {
		return nil, err
	}

	var responseUsers []GetUsersResponse
	for _, user := range users {
		responseUser := GetUsersResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
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

func GetTotalDataSvc(keyword, roleID, active, startDate, endDate, companyID string) (int64, error) {
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

	count, err := GetTotalData(keyword, roleID, active, startTime, endTime, companyID)
	return count, err
}

func DeleteUserByIDSvc(id string) error {
	err := DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
