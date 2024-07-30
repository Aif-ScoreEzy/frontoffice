package auth

import (
	"encoding/json"
	"errors"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/company"
	"front-office/pkg/core/role"
	"front-office/pkg/core/user"
	"front-office/utility/mailjet"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewService(
	repo Repository,
	repoUser user.Repository,
	repoRole role.Repository,
	cfg *config.Config,
) Service {
	return &service{
		Repo:     repo,
		RepoUser: repoUser,
		RepoRole: repoRole,
		Cfg:      cfg,
	}
}

type service struct {
	Repo     Repository
	RepoUser user.Repository
	RepoRole role.Repository
	Cfg      *config.Config
}

type Service interface {
	RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, string, error)
	RegisterMemberSvc(req *user.RegisterMemberRequest, companyId string) (*user.User, string, error)
	VerifyUserTxSvc(userId, token string, req *PasswordResetRequest) (*user.User, error)
	PasswordResetSvc(userId, token string, req *PasswordResetRequest) error
	// LoginSvc(req *UserLoginRequest, user *user.User) (string, string, error)
	ChangePasswordSvc(currentUser *user.User, req *ChangePasswordRequest) (*user.User, error)
	LoginAifCoreService(req *UserLoginRequest, user *user.MstMember) (string, string, error)
	ChangePasswordAifCoreService(req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error)
}

func (svc *service) RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, string, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtVerificationExpiresMinutes)

	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return nil, "", errors.New(constant.InvalidPassword)
	}

	var tierLevel uint
	if req.RoleId != "" {
		result, err := svc.RepoRole.FindOneById(req.RoleId)
		if result == nil {
			return nil, "", errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, "", err
		} else {
			tierLevel = result.TierLevel
		}
	}

	companyId := uuid.NewString()
	dataCompany := &company.Company{
		Id:              companyId,
		CompanyName:     req.CompanyName,
		CompanyAddress:  req.CompanyAddress,
		CompanyPhone:    req.CompanyPhone,
		AgreementNumber: req.AgreementNumber,
		PaymentScheme:   req.PaymentScheme,
		IndustryId:      req.IndustryId,
	}

	userId := uuid.NewString()
	dataUser := &user.User{
		Id:     userId,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Key:    helper.GenerateAPIKey(),
		RoleId: req.RoleId,
	}

	dataUser.Password = user.SetPassword(req.Password)

	token, err := helper.GenerateToken(secret, minutesToExpired, 1, 1, tierLevel)
	if err != nil {
		return nil, "", err
	}

	tokenId := uuid.NewString()
	dataActivationToken := &activationtoken.MstActivationToken{
		Id:     tokenId,
		Token:  token,
		UserId: userId,
	}

	user, err := svc.Repo.CreateAdmin(dataCompany, dataUser, dataActivationToken)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (svc *service) RegisterMemberSvc(req *user.RegisterMemberRequest, companyId string) (*user.User, string, error) {
	userId := uuid.NewString()

	var tierLevel uint
	if req.RoleId != "" {
		result, err := svc.RepoRole.FindOneById(req.RoleId)
		if result == nil {
			return nil, "", errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, "", err
		} else {
			tierLevel = result.TierLevel
		}
	}

	dataUser := &user.User{
		Id:        userId,
		Name:      req.Name,
		Email:     req.Email,
		Key:       helper.GenerateAPIKey(),
		Image:     "default-profile-image.jpg",
		RoleId:    req.RoleId,
		CompanyId: companyId,
	}

	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtActivationExpiresMinutes)

	token, err := helper.GenerateToken(secret, minutesToExpired, 1, 1, tierLevel)
	if err != nil {
		return nil, "", err
	}

	tokenId := uuid.NewString()
	dataToken := &activationtoken.MstActivationToken{
		Id:     tokenId,
		Token:  token,
		UserId: userId,
	}

	user, err := svc.Repo.CreateMember(dataUser, dataToken)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (svc *service) VerifyUserTxSvc(userId, token string, req *PasswordResetRequest) (*user.User, error) {
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

	user, err := svc.Repo.VerifyUserTx(updateUser, userId, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) PasswordResetSvc(userId, token string, req *PasswordResetRequest) error {
	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return errors.New(constant.InvalidPassword)
	}

	if req.Password != req.ConfirmPassword {
		return errors.New(constant.ConfirmPasswordMismatch)
	}

	err := svc.Repo.ResetPassword(userId, token, req)
	if err != nil {
		return err
	}

	return nil
}

// func (svc *service) LoginSvc(req *UserLoginRequest, user *user.User) (string, string, error) {
// 	secret := svc.Cfg.Env.JwtSecretKey

// 	accessTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtExpiresMinutes)
// 	err := bcrypt.CompareHashAndPassword(
// 		[]byte(user.Password),
// 		[]byte(req.Password),
// 	)
// 	if err != nil {
// 		return "", "", errors.New(constant.InvalidEmailOrPassword)
// 	}

// 	accessToken, err := helper.GenerateToken(secret, accessTokenExpiresAt, user.Id, user.CompanyId, user.Role.TierLevel)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	refreshTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtRefreshTokenExpiresMinutes)
// 	refreshToken, err := helper.GenerateRefreshToken(secret, refreshTokenExpiresAt, user.Id, user.CompanyId, user.Role.TierLevel)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	return accessToken, refreshToken, nil
// }

func (svc *service) ChangePasswordSvc(currentUser *user.User, req *ChangePasswordRequest) (*user.User, error) {
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

	data, err := svc.RepoUser.UpdateOneById(updateUser, currentUser)
	if err != nil {
		return nil, err
	}

	err = mailjet.SendConfirmationEmailPasswordChangeSuccess(currentUser.Name, currentUser.Email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (svc *service) LoginAifCoreService(req *UserLoginRequest, user *user.MstMember) (string, string, error) {
	secret := svc.Cfg.Env.JwtSecretKey

	accessTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtExpiresMinutes)
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", "", errors.New(constant.InvalidEmailOrPassword)
	}

	accessToken, err := helper.GenerateToken(secret, accessTokenExpiresAt, user.MemberId, user.CompanyId, user.Role.RoleId)
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtRefreshTokenExpiresMinutes)
	refreshToken, err := helper.GenerateRefreshToken(secret, refreshTokenExpiresAt, user.MemberId, user.CompanyId, user.Role.RoleId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (svc *service) ChangePasswordAifCoreService(req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error) {
	response, err := svc.Repo.ChangePasswordAifCoreService(req)
	if err != nil {
		return nil, err
	}

	var baseResponseSuccess *helper.BaseResponseSuccess
	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		json.Unmarshal(dataBytes, &baseResponseSuccess)
		baseResponseSuccess.StatusCode = response.StatusCode
	}

	return baseResponseSuccess, nil
}
