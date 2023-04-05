package user

import (
	"errors"
	"front-office/helper"
	"front-office/pkg/company"

	"github.com/google/uuid"
)

func RegisterUserSvc(req RegisterUserRequest) (User, error) {
	var user User
	isEmailExist := GetUserByEmailSvc(req.Email)
	if isEmailExist {
		return user, errors.New("Email already exists")
	}

	companyID := uuid.NewString()
	dataCompany := company.Company{
		ID:              companyID,
		CompanyName:     req.CompanyName,
		CompanyAddress:  req.CompanyAddress,
		CompanyPhone:    req.CompanyPhone,
		AgreementNumber: req.AgreementNumber,
		PaymentScheme:   req.PaymentScheme,
		IndustryID:      req.IndustryID,
	}

	userID := uuid.NewString()
	dataUser := User{
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

func GetUserByEmailSvc(email string) bool {
	user := User{
		Email: email,
	}

	result := user.FindOneByEmail()
	return result.ID != ""
}
