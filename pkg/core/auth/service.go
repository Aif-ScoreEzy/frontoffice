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
	VerifyUserTxSvc(userId, token string, req *PasswordResetRequest) (*user.User, error)
	PasswordResetSvc(userId, token string, req *PasswordResetRequest) error
	AddMemberAifCoreService(req *user.RegisterMemberRequest, companyId uint) (*user.RegisterMemberResponse, error)
	LoginMemberAifCoreService(req *UserLoginRequest, user *user.MstMember) (string, string, error)
	ChangePasswordAifCoreService(member *user.FindUserAifCoreResponse, req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error)
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

func (svc *service) PasswordResetSvc(memberId, token string, req *PasswordResetRequest) error {
	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return errors.New(constant.InvalidPassword)
	}

	if req.Password != req.ConfirmPassword {
		return errors.New(constant.ConfirmPasswordMismatch)
	}

	_, err := svc.Repo.PasswordReset(memberId, token, req)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) AddMemberAifCoreService(req *user.RegisterMemberRequest, companyId uint) (*user.RegisterMemberResponse, error) {
	data := &user.RegisterMemberRequest{
		Name:      req.Name,
		Email:     req.Email,
		CompanyId: companyId,
	}

	res, err := svc.RepoUser.AddMemberAifCore(data)
	if err != nil {
		return nil, err
	}

	var baseResponseSuccess *user.RegisterMemberResponse
	if res != nil {
		dataBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		json.Unmarshal(dataBytes, &baseResponseSuccess)
		baseResponseSuccess.StatusCode = res.StatusCode
	}

	return baseResponseSuccess, nil
}

func (svc *service) LoginMemberAifCoreService(req *UserLoginRequest, user *user.MstMember) (string, string, error) {
	secret := svc.Cfg.Env.JwtSecretKey

	accessTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtExpiresMinutes)
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", "", errors.New(constant.InvalidEmailOrPassword)
	}

	accessToken, err := helper.GenerateToken(secret, accessTokenExpiresAt, user.MemberId, user.CompanyId, user.RoleId)
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtRefreshTokenExpiresMinutes)
	refreshToken, err := helper.GenerateToken(secret, refreshTokenExpiresAt, user.MemberId, user.CompanyId, user.RoleId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (svc *service) ChangePasswordAifCoreService(member *user.FindUserAifCoreResponse, req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error) {
	err := bcrypt.CompareHashAndPassword([]byte(member.Data.Password), []byte(req.CurrentPassword))
	if err != nil {
		return nil, errors.New("old_password is wrong")
	}

	isPasswordStrength := helper.ValidatePasswordStrength(req.NewPassword)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	if req.NewPassword != req.ConfirmNewPassword {
		return nil, errors.New(constant.ConfirmNewPasswordMismatch)
	}

	memberIdStr := helper.ConvertUintToString(member.Data.MemberId)
	response, err := svc.Repo.ChangePasswordAifCoreService(memberIdStr, req)
	if err != nil {
		return nil, err
	}

	var baseResponseSuccess *helper.BaseResponseSuccess
	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponseSuccess); err != nil {
			return nil, err
		}
		baseResponseSuccess.StatusCode = response.StatusCode
	}

	return baseResponseSuccess, nil
}
