package activationtoken

import (
	"errors"
	"strconv"

	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/apperror/mapper"
	"strings"
)

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{Repo: repo, Cfg: cfg}
}

type service struct {
	Repo Repository
	Cfg  *config.Config
}

type Service interface {
	CreateActivationToken(memberId, companyId uint, roleId uint) (string, error)
	ValidateActivationToken(authHeader string) (string, uint, error)
	GetActivationToken(token string) (*MstActivationToken, error)
}

func (svc *service) CreateActivationToken(memberId, companyId, roleId uint) (string, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtActivationExpiresMinutes)

	token, err := helper.GenerateToken(secret, minutesToExpired, memberId, companyId, roleId, "")
	if err != nil {
		return "", err
	}

	req := &CreateActivationTokenRequest{
		Token: token,
	}

	memberIdStr := helper.ConvertUintToString(memberId)
	err = svc.Repo.CallCreateActivationTokenAPI(memberIdStr, req)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *service) ValidateActivationToken(authHeader string) (string, uint, error) {
	secret := svc.Cfg.Env.JwtSecretKey

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		return "", 0, errors.New(constant.InvalidActivationLink)
	}

	token := bearerToken[1]

	claims, err := helper.ExtractClaimsFromJWT(token, secret)
	if err != nil {
		return "", 0, errors.New(constant.InvalidActivationLink)
	}

	userId, err := helper.ExtractUserIdFromClaims(claims)
	if err != nil {
		return "", 0, errors.New(constant.InvalidActivationLink)
	}

	return token, userId, nil
}

func (svc *service) GetActivationToken(token string) (*MstActivationToken, error) {
	activationToken, err := svc.Repo.CallGetActivationTokenAPI(token)
	if err != nil {
		return nil, mapper.MapRepoError(err, "failed to get activation token")
	}

	return activationToken, nil
}
