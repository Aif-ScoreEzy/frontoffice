package user

import (
	"errors"
	"front-office/helper"
	"front-office/pkg/company"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserSvc(req RegisterUserRequest) (*User, error) {
	companyID := uuid.NewString()
	dataCompany := &company.Company{
		ID:              companyID,
		CompanyName:     req.CompanyName,
		CompanyAddress:  req.CompanyAddress,
		CompanyPhone:    req.CompanyPhone,
		AgreementNumber: req.AgreementNumber,
		PaymentScheme:   req.PaymentScheme,
		IndustryID:      req.IndustryID,
	}

	userID := uuid.NewString()
	dataUser := &User{
		ID:       userID,
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Key:      helper.GenerateAPIKey(),
		RoleID:   req.RoleID,
	}

	dataUser.SetPassword(req.Password)

	user, err := Create(dataCompany, dataUser)
	if err != nil {
		return user, err
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

func FindUserByUsernameSvc(username string) (*User, error) {
	user, err := FindOneByUsername(username)
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

func LoginSvc(req *UserLoginRequest, user *User) (string, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("password is incorrect")
	}

	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
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
	dataReq := &User{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		CompanyID: req.CompanyID,
		RoleID:    req.RoleID,
	}

	user, err := UpdateOneByID(dataReq, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetAllUsersSvc() ([]*User, error) {
	users, err := FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}
