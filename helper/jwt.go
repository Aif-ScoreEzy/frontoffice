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
	userID, companyID string,
	tierLevel uint,
) (string, error) {
	willExpiredAt := time.Now().Add(time.Duration(minutesToExpired) * time.Minute)

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["company_id"] = companyID
	claims["tier_level"] = tierLevel
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

func ExtractUserIDFromClaims(claims *jwt.MapClaims) (string, error) {
	x, found := (*claims)["user_id"]
	if found {
		if _, ok := x.(string); !ok {
			return "", errors.New("value can't be coerced to string")
		}
	} else {
		return "", errors.New("key doesn't exist")
	}

	return (*claims)["user_id"].(string), nil
}

func ExtractCompanyIDFromClaims(claims *jwt.MapClaims) (string, error) {
	x, found := (*claims)["company_id"]
	if found {
		if _, ok := x.(string); !ok {
			return "", errors.New("value can't be coerced to string")
		}
	} else {
		return "", errors.New("key doesn't exist")
	}

	return (*claims)["company_id"].(string), nil
}

func ExtractTierLevelFromClaims(claims *jwt.MapClaims) (uint, error) {
	x, found := (*claims)["tier_level"]
	if found {
		tierLevelStr := fmt.Sprintf("%v", x)

		tierLevel, err := strconv.ParseUint(tierLevelStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(tierLevel), nil
	} else {
		return 0, errors.New("key doesn't exist")
	}
}
