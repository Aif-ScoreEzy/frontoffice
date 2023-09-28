package user

import (
	"front-office/helper"
	"front-office/utility/mailjet"

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

func GetAllUsersSvc(limit, offset int, companyID string) ([]GetUsersResponse, error) {
	users, err := FindAll(limit, offset, companyID)
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

func DeleteUserByIDSvc(id string) error {
	err := DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
