package auth

import (
	"encoding/json"
	"errors"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"
	"io"
	"strconv"
)

func NewService(
	repo Repository,
	repoUser member.Repository,
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
	RepoUser member.Repository
	RepoRole role.Repository
	Cfg      *config.Config
}

type Service interface {
	// RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, string, error)
	PasswordResetSvc(memberId uint, token string, req *PasswordResetRequest) error
	VerifyMemberAif(memberId uint, req *PasswordResetRequest) (*helper.BaseResponseSuccess, error)
	AddMember(req *member.RegisterMemberRequest, companyId uint) (*member.RegisterMemberResponse, error)
	LoginMember(req *UserLoginRequest) (*aifcoreAuthMemberResponse, error)
	ChangePassword(memberId string, req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error)
	generateTokens(memberId, companyId, roleId uint) (string, string, error)
}

// func (svc *service) RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, string, error) {
// 	secret := svc.Cfg.Env.JwtSecretKey
// 	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtVerificationExpiresMinutes)

// 	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
// 	if !isPasswordStrength {
// 		return nil, "", errors.New(constant.InvalidPassword)
// 	}

// 	var tierLevel uint
// 	if req.RoleId != "" {
// 		result, err := svc.RepoRole.FindOneById(req.RoleId)
// 		if result == nil {
// 			return nil, "", errors.New(constant.DataNotFound)
// 		} else if err != nil {
// 			return nil, "", err
// 		} else {
// 			tierLevel = result.TierLevel
// 		}
// 	}

// 	companyId := uuid.NewString()
// 	dataCompany := &company.Company{
// 		Id:              companyId,
// 		CompanyName:     req.CompanyName,
// 		CompanyAddress:  req.CompanyAddress,
// 		CompanyPhone:    req.CompanyPhone,
// 		AgreementNumber: req.AgreementNumber,
// 		PaymentScheme:   req.PaymentScheme,
// 		IndustryId:      req.IndustryId,
// 	}

// 	memberId := uuid.NewString()
// 	dataUser := &user.User{
// 		Id:     memberId,
// 		Name:   req.Name,
// 		Email:  req.Email,
// 		Phone:  req.Phone,
// 		Key:    helper.GenerateAPIKey(),
// 		RoleId: req.RoleId,
// 	}

// 	dataUser.Password = user.SetPassword(req.Password)

// 	token, err := helper.GenerateToken(secret, minutesToExpired, 1, 1, tierLevel)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	tokenId := uuid.NewString()
// 	dataActivationToken := &activationtoken.MstActivationToken{
// 		Id:       tokenId,
// 		Token:    token,
// 		MemberId: memberId,
// 	}

// 	user, err := svc.Repo.CreateAdmin(dataCompany, dataUser, dataActivationToken)
// 	if err != nil {
// 		return user, "", err
// 	}

// 	return user, token, nil
// }

func (svc *service) VerifyMemberAif(memberId uint, req *PasswordResetRequest) (*helper.BaseResponseSuccess, error) {
	isPasswordStrength := helper.ValidatePasswordStrength(req.Password)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New(constant.ConfirmPasswordMismatch)
	}

	response, err := svc.Repo.VerifyMemberAif(req, memberId)
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

func (svc *service) PasswordResetSvc(memberId uint, token string, req *PasswordResetRequest) error {
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

func (svc *service) AddMember(req *member.RegisterMemberRequest, companyId uint) (*member.RegisterMemberResponse, error) {
	res, err := svc.RepoUser.AddMember(req)
	if err != nil {
		return nil, err
	}

	var baseResponseSuccess *member.RegisterMemberResponse
	if res != nil {
		dataBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		json.Unmarshal(dataBytes, &baseResponseSuccess)
		baseResponseSuccess.StatusCode = res.StatusCode
	}

	return baseResponseSuccess, nil
}

func (svc *service) LoginMember(req *UserLoginRequest) (*aifcoreAuthMemberResponse, error) {
	res, err := svc.Repo.AuthMemberAifCore(req)
	if err != nil {
		return nil, err
	}

	var baseResponse *aifcoreAuthMemberResponse
	if res != nil {
		dataBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		json.Unmarshal(dataBytes, &baseResponse)
		baseResponse.StatusCode = res.StatusCode
	}

	return baseResponse, nil
}

func (svc *service) ChangePassword(memberId string, req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error) {
	isPasswordStrength := helper.ValidatePasswordStrength(req.NewPassword)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	if req.NewPassword != req.ConfirmNewPassword {
		return nil, errors.New(constant.ConfirmNewPasswordMismatch)
	}

	response, err := svc.Repo.ChangePasswordAifCore(memberId, req)
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

func (svc *service) generateTokens(memberId, companyId, roleId uint) (string, string, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	accessTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtExpiresMinutes)
	accessToken, err := helper.GenerateToken(secret, accessTokenExpiresAt, memberId, companyId, roleId)
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiresAt, _ := strconv.Atoi(svc.Cfg.Env.JwtRefreshTokenExpiresMinutes)
	refreshToken, err := helper.GenerateToken(secret, refreshTokenExpiresAt, memberId, companyId, roleId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
