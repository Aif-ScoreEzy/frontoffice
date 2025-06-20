package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/internal/apperror/mapper"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"
	"io"
	"strconv"

	"github.com/rs/zerolog/log"
)

func NewService(
	cfg *config.Config,
	repo Repository,
	memberRepo member.Repository,
	roleRepo role.Repository,
	operationRepo operation.Repository,
) Service {
	return &service{
		cfg,
		repo,
		memberRepo,
		roleRepo,
		operationRepo,
	}
}

type service struct {
	cfg           *config.Config
	repo          Repository
	memberRepo    member.Repository
	roleRepo      role.Repository
	operationRepo operation.Repository
}

type Service interface {
	// RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, string, error)
	PasswordResetSvc(memberId uint, token string, req *PasswordResetRequest) error
	VerifyMemberAif(memberId uint, req *PasswordResetRequest) (*helper.BaseResponseSuccess, error)
	AddMember(req *member.RegisterMemberRequest, companyId uint) (*member.RegisterMemberResponse, error)
	LoginMember(loginReq *userLoginRequest) (accessToken, refreshToken string, loginResp *loginResponse, err error)
	ChangePassword(memberId string, req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error)
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

	response, err := svc.repo.VerifyMemberAif(req, memberId)
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

	_, err := svc.repo.PasswordReset(memberId, token, req)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) AddMember(req *member.RegisterMemberRequest, companyId uint) (*member.RegisterMemberResponse, error) {
	res, err := svc.memberRepo.AddMember(req)
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

func (svc *service) LoginMember(req *userLoginRequest) (accessToken, refreshToken string, loginResp *loginResponse, err error) {
	user, err := svc.repo.AuthMemberAifCore(req)
	if err != nil {
		var apiErr *apperror.ExternalAPIError
		if errors.As(err, &apiErr) {
			return "", "", nil, mapper.MapAuthError(apiErr)
		}

		return "", "", nil, apperror.Internal("auth failed", err)
	}

	accessToken, err = svc.generateToken(user, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtExpiresMinutes)
	if err != nil {
		return "", "", nil, apperror.Internal("generate access token failed", err)
	}

	refreshToken, err = svc.generateToken(user, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtRefreshTokenExpiresMinutes)
	if err != nil {
		return "", "", nil, apperror.Internal("generate refresh token failed", err)
	}

	_, err = svc.operationRepo.AddLogOperation(&operation.AddLogRequest{
		MemberId:  user.MemberId,
		CompanyId: user.CompanyId,
		Action:    constant.EventSignIn,
	})
	if err != nil {
		log.Warn().Err(err).Msg("failed to log sign-in event")
	}

	loginResp = &loginResponse{
		Id:                 user.MemberId,
		Name:               user.Name,
		Email:              user.Email,
		CompanyId:          user.CompanyId,
		CompanyName:        user.CompanyName,
		TierLevel:          user.RoleId,
		Image:              user.Image,
		SubscriberProducts: user.SubscriberProducts,
	}

	return accessToken, refreshToken, loginResp, nil
}

func (svc *service) ChangePassword(memberId string, req *ChangePasswordRequest) (*helper.BaseResponseSuccess, error) {
	isPasswordStrength := helper.ValidatePasswordStrength(req.NewPassword)
	if !isPasswordStrength {
		return nil, errors.New(constant.InvalidPassword)
	}

	if req.NewPassword != req.ConfirmNewPassword {
		return nil, errors.New(constant.ConfirmNewPasswordMismatch)
	}

	response, err := svc.repo.ChangePasswordAifCore(memberId, req)
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

func (svc *service) generateToken(user *loginResponseData, secret, minutesStr string) (string, error) {
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return "", fmt.Errorf("invalid duration: %w", err)
	}

	return helper.GenerateToken(secret, minutes, user.MemberId, user.CompanyId, user.RoleId, user.ApiKey)
}
