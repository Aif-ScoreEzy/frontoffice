package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(
	secret string,
	minutesToExpired int,
	userID string,
	tierLevel uint,
) (string, error) {
	willExpiredAt := time.Now().Add(time.Duration(minutesToExpired) * time.Minute)

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
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
