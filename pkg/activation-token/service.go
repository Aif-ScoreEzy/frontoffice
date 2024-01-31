package activation_token

import (
	"errors"

	"front-office/common/constant"
	"front-office/helper"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func CreateActivationTokenSvc(userID, companyID string, tierLevel uint) (string, *ActivationToken, error) {
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

	activationToken, err = CreateActivationToken(activationToken)
	if err != nil {
		return "", nil, err
	}

	return token, activationToken, nil
}

func ValidateActivationToken(authHeader string) (string, string, error) {
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

func FindActivationTokenByTokenSvc(token string) (*ActivationToken, error) {
	result, err := FindOneActivationTokenBytoken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindActivationTokenByUserIDSvc(userID string) (*ActivationToken, error) {
	result, err := FindOneActivationTokenByUserID(userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
