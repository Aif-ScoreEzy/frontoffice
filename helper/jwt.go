package helper

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(
	userID string,
) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_MINUTES"))
	willExpiredAt := time.Now().Add(time.Duration(minutesToExpired) * time.Minute)

	claims := jwt.MapClaims{}
	claims["userID"] = userID
	claims["exp"] = willExpiredAt.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

type TokenMetadata struct {
	Expires int64
	UserID  string
}

func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		expires := int64(claims["exp"].(float64))
		userID := claims["userID"].(string)

		return &TokenMetadata{
			Expires: expires,
			UserID:  userID,
		}, nil
	}

	return nil, err
}

func extractToken(c *fiber.Ctx) string {
	beareraToken := c.Get("Authorization")

	onlyToken := strings.Split(beareraToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
