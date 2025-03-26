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
	CreateActivationToken(userId, companyId uint, roleId uint) (string, *AifResponse, error)
	ValidateActivationToken(authHeader string) (string, uint, error)
	FindActivationTokenByToken(token string) (*AifResponse, error)
}

func (svc *service) CreateActivationToken(userId, companyId, roleId uint) (string, *AifResponse, error) {
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
	response, err := svc.Repo.CreateActivationTokenAifCore(req, userIdStr)
	if err != nil {
		return "", nil, err
	}

	var baseResponse *AifResponse
	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponse); err != nil {
			return "", nil, err
		}
	}

	return token, baseResponse, nil
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

func (svc *service) FindActivationTokenByToken(token string) (*AifResponse, error) {
	response, err := svc.Repo.FindOneActivationTokenBytoken(token)
	if err != nil {
		return nil, err
	}

	var baseResponse *AifResponse
	if response != nil {
		dataBytes, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponse); err != nil {
			return nil, err
		}
	}

	return baseResponse, nil
}
