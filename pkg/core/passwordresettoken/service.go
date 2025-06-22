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
	FindPasswordResetTokenByTokenSvc(token string) (*FindTokenResponse, error)
	CreatePasswordResetToken(userId, companyId, roleId uint) (string, error)
	DeletePasswordResetToken(id uint) (*helper.BaseResponseSuccess, error)
}

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

func (svc *service) CreatePasswordResetToken(userId, companyId, roleId uint) (string, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtActivationExpiresMinutes)

	token, err := helper.GenerateToken(secret, minutesToExpired, userId, companyId, roleId, "")
	if err != nil {
		return "", err
	}

	req := &CreatePasswordResetTokenRequest{
		Token: token,
	}

	userIdStr := helper.ConvertUintToString(userId)
	err = svc.Repo.CallCreatePasswordResetToken(userIdStr, req)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *service) DeletePasswordResetToken(id uint) (*helper.BaseResponseSuccess, error) {
	idStr := strconv.Itoa(int(id))
	response, err := svc.Repo.DeletePasswordResetToken(idStr)
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
