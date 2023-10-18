package auth

import (
	"errors"
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/company"
	"front-office/pkg/role"
	"front-office/pkg/user"
	"front-office/utility/mailjet"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, error) {
	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	var tierLevel uint
	if req.RoleID != "" {
		result, err := role.FindRoleByIDSvc(req.RoleID)
		if result == nil {
			return nil, errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, err
		} else {
			tierLevel = result.TierLevel
		}
	}

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
	dataUser := &user.User{
		ID:     userID,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Key:    helper.GenerateAPIKey(),
		RoleID: req.RoleID,
	}

	dataUser.Password = user.SetPassword(req.Password)

	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_ACTIVATION_EXPIRES_MINUTES"))

	token, err := helper.GenerateToken(secret, minutesToExpired, userID, companyID, tierLevel)
	if err != nil {
		return nil, err
	}

	tokenID := uuid.NewString()
	dataActivationToken := &user.ActivationToken{
		ID:     tokenID,
		Token:  token,
		UserID: userID,
	}

	user, err := CreateAdmin(dataCompany, dataUser, dataActivationToken)
	if err != nil {
		return user, err
	}

	return user, nil
}

func SendEmailVerificationSvc(req *SendEmailVerificationRequest, user *user.User) error {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_VERIFICATION_EXPIRES_MINUTES"))
	baseURL := os.Getenv("FRONTEND_BASE_URL")

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID, user.CompanyID, user.Role.TierLevel)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%s/verify/%s", baseURL, token),
	}

	err = mailjet.CreateMailjet(req.Email, 5075167, variables)
	if err != nil {
		return err
	}

	return nil
}

func VerifyUserTxSvc(userID, token string, req *PasswordResetRequest) (*user.User, error) {
	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New(constant.ConfirmPasswordMismatch)
	}

	updateUser := map[string]interface{}{}

	updateUser["password"] = user.SetPassword(req.Password)
	updateUser["status"] = "active"
	updateUser["active"] = true
	updateUser["is_verified"] = true
	updateUser["updated_at"] = time.Now()

	user, err := user.VerifyUserTx(updateUser, userID, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_RESET_PASSWORD_EXPIRES_MINUTES"))

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID, user.CompanyID, user.Role.TierLevel)
	if err != nil {
		return "", nil, err
	}

	tokenID := uuid.NewString()
	passwordResetToken := &PasswordResetToken{
		ID:     tokenID,
		Token:  token,
		UserID: user.ID,
	}

	passwordResetToken, err = CreatePasswordResetToken(passwordResetToken)
	if err != nil {
		return "", nil, err
	}

	return token, passwordResetToken, nil
}

func SendEmailPasswordResetSvc(email, token string) error {
	baseURL := os.Getenv("FRONTEND_BASE_URL")

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%s/verification?key=%s", baseURL, token),
	}

	err := mailjet.CreateMailjet(email, 5085661, variables)
	if err != nil {
		return err
	}

	return nil
}

func FindPasswordResetTokenByTokenSvc(token string) (*PasswordResetToken, error) {
	result, err := FindOnePasswordTokenByToken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PasswordResetSvc(userID, token string, req *PasswordResetRequest) error {
	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return errors.New(constant.InvalidPassword)
	}

	if req.Password != req.ConfirmPassword {
		return errors.New(constant.ConfirmPasswordMismatch)
	}

	err := ResetPassword(userID, token, req)
	if err != nil {
		return err
	}

	return nil
}

func LoginSvc(req *UserLoginRequest, user *user.User) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_MINUTES"))
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New(constant.InvalidEmailOrPassword)
	}

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID, user.CompanyID, user.Role.TierLevel)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ChangePasswordSvc(currentUser *user.User, req *ChangePasswordRequest) (*user.User, error) {
	updateUser := map[string]interface{}{}

	err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.CurrentPassword))
	if err != nil {
		return nil, errors.New(constant.IncorrectPassword)
	}

	isPasswordStrength := helper.ValidatePasswordStrength(req.NewPassword)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	if req.NewPassword != req.ConfirmNewPassword {
		return nil, errors.New(constant.ConfirmNewPasswordMismatch)
	}

	updateUser["password"] = user.SetPassword(req.NewPassword)
	updateUser["updated_at"] = time.Now()

	data, err := user.UpdateOneByID(updateUser, currentUser)
	if err != nil {
		return nil, err
	}

	variables := map[string]interface{}{
		"username": currentUser.Name,
	}

	err = mailjet.CreateMailjet(currentUser.Email, 5097353, variables)
	if err != nil {
		return nil, err
	}

	return data, nil
}
