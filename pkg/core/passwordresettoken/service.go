package passwordresettoken

import (
	"encoding/json"
	"front-office/app/config"
	"front-office/helper"
	"io"
	"strconv"
)

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{Repo: repo, Cfg: cfg}
}

type service struct {
	Repo Repository
	Cfg  *config.Config
}

type Service interface {
	// CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error)
	FindPasswordResetTokenByTokenSvc(token string) (*FindTokenResponse, error)
	CreatePasswordResetTokenAifCore(userId, companyId, roleId uint) (string, error)
}

// func (svc *service) CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error) {
// 	secret := svc.Cfg.Env.JwtSecretKey
// 	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtResetPasswordExpiresMinutes)

// 	token, err := helper.GenerateToken(secret, minutesToExpired, user.Id, user.CompanyId, user.Role.TierLevel)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	tokenId := uuid.NewString()
// 	passwordResetToken := &PasswordResetToken{
// 		Id:     tokenId,
// 		Token:  token,
// 		UserId: user.Id,
// 	}

// 	passwordResetToken, err = svc.Repo.CreatePasswordResetToken(passwordResetToken)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	return token, passwordResetToken, nil
// }

func (svc *service) FindPasswordResetTokenByTokenSvc(token string) (*FindTokenResponse, error) {
	response, err := svc.Repo.FindOnePasswordResetTokenByToken(token)
	if err != nil {
		return nil, err
	}

	var baseResponseSuccess *FindTokenResponse
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

func (svc *service) CreatePasswordResetTokenAifCore(userId, companyId, roleId uint) (string, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtActivationExpiresMinutes)

	token, err := helper.GenerateToken(secret, minutesToExpired, userId, companyId, roleId)
	if err != nil {
		return "", err
	}

	req := &CreatePasswordResetTokenRequest{
		Token: token,
	}

	userIdStr := helper.ConvertUintToString(userId)
	_, err = svc.Repo.CreatePasswordResetTokenAifCore(req, userIdStr)
	if err != nil {
		return "", err
	}

	return token, nil
}
