package user

import (
	"errors"
	"front-office/helper"
	"front-office/pkg/company"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserSvc(req RegisterUserRequest) (User, error) {
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

func IsEmailExistSvc(email string) bool {
	user := User{
		Email: email,
	}

	result := user.FindOneByEmail()
	return result.ID != ""
}

func IsUsernameExistSvc(username string) (bool, User) {
	user := User{
		Username: username,
	}

	result := user.FindOneByUsername()

	return result.ID != "", result
}

func LoginSvc(req UserLoginRequest, user User) (string, error) {
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
