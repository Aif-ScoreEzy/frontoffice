package activationtoken

import (
	"encoding/json"
	"errors"
	"io"
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
	CreateActivationTokenSvc(userId, companyId uint, tierLevel uint) (string, *CreateActivationTokenResponse, error)
	ValidateActivationToken(authHeader string) (string, uint, error)
	FindActivationTokenByTokenSvc(token string) (*MstActivationToken, error)
	FindActivationTokenByUserIdSvc(userId string) (*MstActivationToken, error)
}

func (svc *service) CreateActivationTokenSvc(userId, companyId, roleId uint) (string, *CreateActivationTokenResponse, error) {
	secret := svc.Cfg.Env.JwtSecretKey
	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtActivationExpiresMinutes)

	token, err := helper.GenerateToken(secret, minutesToExpired, userId, companyId, roleId)
	if err != nil {
		return "", nil, err
	}

	req := &CreateActivationTokenRequest{
		Token: token,
	}

	userIdStr := helper.ConvertUintToString(userId)
	res, err := svc.Repo.CreateActivationTokenAifCore(req, userIdStr)
	if err != nil {
		return "", nil, err
	}

	var baseResponseSuccess *CreateActivationTokenResponse
	if res != nil {
		dataBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		json.Unmarshal(dataBytes, &baseResponseSuccess)
		baseResponseSuccess.StatusCode = res.StatusCode
	}

	return token, baseResponseSuccess, nil
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
