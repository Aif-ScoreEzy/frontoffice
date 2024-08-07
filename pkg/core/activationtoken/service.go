package activationtoken

import (
	"errors"
	"strconv"

	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
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
	CreateActivationTokenAifCore(userId, companyId uint, roleId uint) (string, error)
	ValidateActivationToken(authHeader string) (string, uint, error)
	FindActivationTokenByTokenSvc(token string) (*MstActivationToken, error)
	FindActivationTokenByUserIdSvc(userId string) (*MstActivationToken, error)
}

func (svc *service) CreateActivationTokenAifCore(userId, companyId, roleId uint) (string, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtActivationExpiresMinutes)

	token, err := helper.GenerateToken(secret, minutesToExpired, userId, companyId, roleId)
	if err != nil {
		return "", err
	}

	req := &CreateActivationTokenRequest{
		Token: token,
	}

	userIdStr := helper.ConvertUintToString(userId)
	_, err = svc.Repo.CreateActivationTokenAifCore(req, userIdStr)
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

func (svc *service) FindActivationTokenByTokenSvc(token string) (*MstActivationToken, error) {
	result, err := svc.Repo.FindOneActivationTokenBytoken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) FindActivationTokenByUserIdSvc(userId string) (*MstActivationToken, error) {
	result, err := svc.Repo.FindOneActivationTokenByUserId(userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
