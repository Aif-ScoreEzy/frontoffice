package user

import (
	"errors"
	"fmt"
	"front-office/helper"
	"front-office/pkg/company"
	"front-office/utility/mailjet"
	"os"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserSvc(req *RegisterUserRequest) (*User, error) {
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
		ID:         userID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Key:        helper.GenerateAPIKey(),
		IsVerified: false,
		RoleID:     req.RoleID,
	}

	dataUser.Password = SetPassword(req.Password)

	user, err := Create(dataCompany, dataUser)
	if err != nil {
		return user, err
	}

	return user, nil
}

func RegisterMemberSvc(req *RegisterMemberRequest, loggedUser *User) (*User, error) {
	userID := uuid.NewString()
	dataUser := &User{
		ID:         userID,
		Name:       req.Name,
		Email:      req.Email,
		Key:        helper.GenerateAPIKey(),
		RoleID:     req.RoleID,
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

	err = mailjet.CreateMailjet(req.Email, "Successful Registration", 5082139, variables)
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
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_MINUTES"))
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("password is incorrect")
	}

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func SendEmailVerificationSvc(req *SendEmailVerificationRequest, user *User) error {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_EMAIL_VERIFICATION_EXPIRES_MINUTES"))
	baseURL := os.Getenv("BASE_URL")

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%s/verify/%s", baseURL, token),
	}

	err = mailjet.CreateMailjet(req.Email, "Email Verification", 5075167, variables)
	if err != nil {
		return err
	}

	return nil
}

func VerifyUserSvc(userID string) (*User, error) {
	req := &User{
		IsVerified: true,
	}

	user, err := UpdateOneByID(req, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func SendEmailPasswordResetSvc(req *RequestPasswordResetRequest, user *User) error {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_RESET_PASSWORD_EXPIRES_MINUTES"))
	baseURL := os.Getenv("BASE_URL")

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%s/password-reset/%s", baseURL, token),
	}

	err = mailjet.CreateMailjet(req.Email, "Reset Password", 5085661, variables)
	if err != nil {
		return err
	}

	return nil
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
		Name: req.Name,
		// Username:  req.Username,
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
