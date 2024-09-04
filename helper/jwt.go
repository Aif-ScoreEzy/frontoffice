package helper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	keyEmpty = "key doesn't exist"
)

func GenerateToken(
	secret string,
	minutesToExpired int,
	userId, companyId, roleId uint,
) (string, error) {
	willExpiredAt := time.Now().Add(time.Duration(minutesToExpired) * time.Minute)

	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["company_id"] = companyId
	claims["role_id"] = roleId
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

func ExtractUserIdFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["user_id"]
	if found {
		userIdStr := fmt.Sprintf("%v", x)

		roleId, err := strconv.ParseUint(userIdStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(roleId), nil
	} else {
		return 0, errors.New(keyEmpty)
	}
}

func ExtractCompanyIdFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["company_id"]
	if found {
		companyIdStr := fmt.Sprintf("%v", x)

		roleId, err := strconv.ParseUint(companyIdStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(roleId), nil
	} else {
		return 0, errors.New(keyEmpty)
	}
}

func ExtractRoleIdFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["role_id"]
	if found {
		roleIdStr := fmt.Sprintf("%v", x)

		roleId, err := strconv.ParseUint(roleIdStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(roleId), nil
	} else {
		return 0, errors.New(keyEmpty)
	}
}
