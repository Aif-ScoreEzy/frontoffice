package helper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(
	secret string,
	minutesToExpired int,
	userID, companyID uint,
	roleID uint,
) (string, error) {
	willExpiredAt := time.Now().Add(time.Duration(minutesToExpired) * time.Minute)

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["company_id"] = companyID
	claims["role_id"] = roleID
	claims["exp"] = willExpiredAt.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GenerateRefreshToken(
	secret string,
	minutesToExpired int,
	userID, companyID uint,
	roleID uint,
) (string, error) {
	willExpiredAt := time.Now().Add(time.Duration(minutesToExpired) * time.Minute)

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["company_id"] = companyID
	claims["role_id"] = roleID
	claims["exp"] = willExpiredAt.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ExtractClaimsFromJWT(token, secret string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(requestToken *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func ExtractUserIDFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["user_id"]
	if found {
		userIDStr := fmt.Sprintf("%v", x)

		roleID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(roleID), nil
	} else {
		return 0, errors.New("key doesn't exist")
	}
}

func ExtractCompanyIDFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["company_id"]
	if found {
		companyIDStr := fmt.Sprintf("%v", x)

		roleID, err := strconv.ParseUint(companyIDStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(roleID), nil
	} else {
		return 0, errors.New("key doesn't exist")
	}
}

func ExtractRoleIDFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["role_id"]
	if found {
		roleIDStr := fmt.Sprintf("%v", x)

		roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(roleID), nil
	} else {
		return 0, errors.New("key doesn't exist")
	}
}
