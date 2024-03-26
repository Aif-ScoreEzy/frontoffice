package activationtoken

import (
	"errors"

	"front-office/common/constant"
	"front-office/helper"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateActivationTokenSvc(userID, companyID string, tierLevel uint) (string, *ActivationToken, error)
	ValidateActivationToken(authHeader string) (string, string, error)
	FindActivationTokenByTokenSvc(token string) (*ActivationToken, error)
	FindActivationTokenByUserIDSvc(userID string) (*ActivationToken, error)
}

func (svc *service) CreateActivationTokenSvc(userID, companyID string, tierLevel uint) (string, *ActivationToken, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_ACTIVATION_EXPIRES_MINUTES"))

	token, err := helper.GenerateToken(secret, minutesToExpired, userID, companyID, tierLevel)
	if err != nil {
		return "", nil, err
	}

	tokenID := uuid.NewString()
	activationToken := &ActivationToken{
		ID:     tokenID,
		Token:  token,
		UserID: userID,
	}

	activationToken, err = svc.Repo.CreateActivationToken(activationToken)
	if err != nil {
		return "", nil, err
	}

	return token, activationToken, nil
}

func (svc *service) ValidateActivationToken(authHeader string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		return "", "", errors.New(constant.InvalidActivationLink)
	}

	token := bearerToken[1]

	claims, err := helper.ExtractClaimsFromJWT(token, secret)
	if err != nil {
		return "", "", errors.New(constant.InvalidActivationLink)
	}

	userID, err := helper.ExtractUserIDFromClaims(claims)
	if err != nil {
		return "", "", errors.New(constant.InvalidActivationLink)
	}

	return token, userID, nil
}

func (svc *service) FindActivationTokenByTokenSvc(token string) (*ActivationToken, error) {
	result, err := svc.Repo.FindOneActivationTokenBytoken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) FindActivationTokenByUserIDSvc(userID string) (*ActivationToken, error) {
	result, err := svc.Repo.FindOneActivationTokenByUserID(userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
