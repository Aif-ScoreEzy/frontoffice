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
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func NewService(
	cfg *config.Config,
	repo Repository,
	memberRepo member.Repository,
	roleRepo role.Repository,
	operationRepo operation.Repository,
	activationRepo activationtoken.Repository,
) Service {
	return &service{
		cfg,
		repo,
		memberRepo,
		roleRepo,
		operationRepo,
		activationRepo,
	}
}

type service struct {
	cfg            *config.Config
	repo           Repository
	memberRepo     member.Repository
	roleRepo       role.Repository
	operationRepo  operation.Repository
	activationRepo activationtoken.Repository
}

type Service interface {
	// RegisterAdminSvc(req *RegisterAdminRequest) (*user.User, string, error)
	LoginMember(loginReq *userLoginRequest) (accessToken, refreshToken string, loginResp *loginResponse, err error)
	RefreshAccessToken(memberId, companyId, tierLevel uint, apiKey string) (string, error)
	Logout(memberId, companyId uint) error
	AddMember(currentUserId uint, req *member.RegisterMemberRequest) error
	SendEmailActivation(email string) error
	PasswordResetSvc(memberId uint, token string, req *PasswordResetRequest) error
	VerifyMember(token string, req *PasswordResetRequest) error
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

func (svc *service) VerifyMember(token string, req *PasswordResetRequest) error {
	activationData, err := svc.activationRepo.CallGetActivationTokenAPI(token)
	if err != nil {
		return mapper.MapRepoError(err, "failed to retrieve activation token")
	}

	memberId := fmt.Sprintf("%d", activationData.MemberId)

	memberData, err := svc.memberRepo.CallGetMemberAPI(&member.FindUserQuery{
		Id: memberId,
	})
	if err != nil {
		return mapper.MapRepoError(err, "failed to fetch member")
	}

	if memberData.IsVerified && memberData.Active {
		return apperror.BadRequest(constant.AlreadyVerified)
	}

	minutesToExpired, err := strconv.Atoi(svc.cfg.Env.JwtActivationExpiresMinutes)
	if err != nil {
		return apperror.Internal("invalid activation expiry config", err)
	}

	elapsedMinutes := time.Since(activationData.CreatedAt).Minutes()
	if elapsedMinutes > float64(minutesToExpired) {
		updateFields := map[string]interface{}{
			"mail_status": mailStatusResend,
			"updated_at":  time.Now(),
		}

		err := svc.memberRepo.CallUpdateMemberAPI(memberId, updateFields)
		if err != nil {
			return mapper.MapRepoError(err, "failed to update member after token expired")
		}

		return apperror.Forbidden(constant.ActivationTokenExpired)
	}

	if !helper.ValidatePasswordStrength(req.Password) {
		return apperror.BadRequest(constant.InvalidPassword)
	}

	if req.Password != req.ConfirmPassword {
		return apperror.BadRequest(constant.ConfirmPasswordMismatch)
	}

	if err := svc.repo.CallVerifyMemberAPI(activationData.MemberId, req); err != nil {
		return mapper.MapRepoError(err, "failed to verify member")
	}

	return nil
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

func (svc *service) AddMember(currentUserId uint, req *member.RegisterMemberRequest) error {
	memberResp, err := svc.memberRepo.CallAddMemberAPI(req)
	if err != nil {
		return mapper.MapRepoError(err, "failed to register member")
	}

	tokenPayload := &tokenPayload{
		MemberId:  memberResp.MemberId,
		CompanyId: req.CompanyId,
		RoleId:    req.RoleId,
	}
	activationToken, err := svc.generateToken(tokenPayload, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtActivationExpiresMinutes)
	if err != nil {
		return apperror.Internal("generate activation token failed", err)
	}

	memberIdStr := helper.ConvertUintToString(memberResp.MemberId)

	err = svc.activationRepo.CallCreateActivationTokenAPI(memberIdStr, &activationtoken.CreateActivationTokenRequest{
		Token: activationToken,
	})
	if err != nil {
		return mapper.MapRepoError(err, "failed to create activation")
	}

	err = mailjet.SendEmailActivation(req.Email, activationToken)
	if err != nil {
		updateFields := map[string]interface{}{
			"mail_status": mailStatusResend,
			"updated_at":  time.Now(),
		}

		err := svc.memberRepo.CallUpdateMemberAPI(memberIdStr, updateFields)
		if err != nil {
			return mapper.MapRepoError(err, "failed to update member after email failure")
		}

		return apperror.Internal("failed to send activation email", err)
	}

	err = svc.operationRepo.AddLogOperation(&operation.AddLogRequest{
		MemberId:  currentUserId,
		CompanyId: req.CompanyId,
		Action:    constant.EventRegisterMember,
	})
	if err != nil {
		log.Warn().Err(err).Msg("failed to log register member event")
	}

	return nil
}

func (svc *service) SendEmailActivation(email string) error {
	memberData, err := svc.memberRepo.CallGetMemberAPI(&member.FindUserQuery{
		Email: email,
	})
	if err != nil {
		return mapper.MapRepoError(err, "failed to fetch member")
	}

	if memberData.IsVerified {
		return apperror.Conflict(constant.AlreadyVerified)
	}

	tokenPayload := &tokenPayload{
		MemberId:  memberData.MemberId,
		CompanyId: memberData.CompanyId,
		RoleId:    memberData.RoleId,
	}
	token, err := svc.generateToken(tokenPayload, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtActivationExpiresMinutes)
	if err != nil {
		return apperror.Internal("generate activation token failed", err)
	}

	memberIdStr := helper.ConvertUintToString(memberData.MemberId)

	if err := svc.activationRepo.CallCreateActivationTokenAPI(memberIdStr, &activationtoken.CreateActivationTokenRequest{
		Token: token,
	}); err != nil {
		return mapper.MapRepoError(err, "failed to create activation")
	}

	if err := mailjet.SendEmailActivation(email, token); err != nil {
		return apperror.Internal("failed to send activation email", err)
	}

	updateFields := map[string]interface{}{
		"mail_status": mailStatusPending,
		"updated_at":  time.Now(),
	}

	if err := svc.memberRepo.CallUpdateMemberAPI(memberIdStr, updateFields); err != nil {
		return mapper.MapRepoError(err, "failed to update member")
	}

	return nil
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

	tokenPayload := &tokenPayload{
		MemberId:  user.MemberId,
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
		ApiKey:    user.ApiKey,
	}
	accessToken, err = svc.generateToken(tokenPayload, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtExpiresMinutes)
	if err != nil {
		return "", "", nil, apperror.Internal("generate access token failed", err)
	}

	refreshToken, err = svc.generateToken(tokenPayload, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtRefreshTokenExpiresMinutes)
	if err != nil {
		return "", "", nil, apperror.Internal("generate refresh token failed", err)
	}

	if err := svc.operationRepo.AddLogOperation(&operation.AddLogRequest{
		MemberId:  user.MemberId,
		CompanyId: user.CompanyId,
		Action:    constant.EventSignIn,
	}); err != nil {
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

func (svc *service) RefreshAccessToken(memberId, companyId, roleId uint, apiKey string) (string, error) {
	tokenPayload := &tokenPayload{
		MemberId:  memberId,
		CompanyId: companyId,
		RoleId:    roleId,
		ApiKey:    apiKey,
	}

	accessToken, err := svc.generateToken(tokenPayload, svc.cfg.Env.JwtSecretKey, svc.cfg.Env.JwtExpiresMinutes)
	if err != nil {
		return "", apperror.Internal("generate access token failed", err)
	}

	return accessToken, nil
}

func (svc *service) Logout(memberId, companyId uint) error {
	if err := svc.operationRepo.AddLogOperation(&operation.AddLogRequest{
		MemberId:  memberId,
		CompanyId: companyId,
		Action:    constant.EventSignOut,
	}); err != nil {
		log.Warn().Err(err).Msg("failed to log sign-out event")
	}

	return nil
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

func (svc *service) generateToken(payload *tokenPayload, secret, minutesStr string) (string, error) {
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return "", fmt.Errorf("invalid duration: %w", err)
	}

	return helper.GenerateToken(secret, minutes, payload.MemberId, payload.CompanyId, payload.RoleId, payload.ApiKey)
}
