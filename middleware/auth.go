package middleware

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v3"
)

func Auth() func(c *fiber.Ctx) error {
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp := helper.ResponseFailed(err.Error())

	return c.Status(fiber.StatusUnauthorized).JSON(resp)
}

func SetHeaderAuth(c *fiber.Ctx) error {
	token := c.Params("token")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return c.Next()
}

func GetUserIDFromJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := os.Getenv("JWT_SECRET_KEY")
		authHeader := c.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			resp := helper.ResponseFailed("Invalid token")

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		token := bearerToken[1]

		claims, err := helper.ExtractClaimsFromJWT(token, secret)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		userID, err := helper.ExtractUserIDFromClaims(claims)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		companyID, err := helper.ExtractCompanyIDFromClaims(claims)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		c.Locals("userID", userID)
		c.Locals("companyID", companyID)

		return c.Next()
	}
}

func AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := os.Getenv("JWT_SECRET_KEY")
		authHeader := c.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			resp := helper.ResponseFailed("Invalid token")

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		token := bearerToken[1]

		claims, err := helper.ExtractClaimsFromJWT(token, secret)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		tierLevel, err := helper.ExtractTierLevelFromClaims(claims)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}
		if tierLevel == 2 {
			resp := helper.ResponseFailed(constant.RequestProhibited)

			return c.Status(fiber.StatusUnauthorized).JSON(resp)
		}

		return c.Next()
	}
}
