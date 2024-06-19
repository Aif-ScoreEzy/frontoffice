package middleware

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Auth() func(c *fiber.Ctx) error {
	config := jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ErrorHandler: jwtError,
		TokenLookup:  "cookie:access_token",
	}

	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	resp := helper.ResponseFailed(err.Error())
	return c.Status(fiber.StatusUnauthorized).JSON(resp)
}

func SetHeaderAuth(c *fiber.Ctx) error {
	token := c.Params("token")
	c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return c.Next()
}

func GetPayloadFromJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := os.Getenv("JWT_SECRET_KEY")
		token := c.Cookies("access_token")

		claims, err := helper.ExtractClaimsFromJWT(token, secret)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusUnauthorized).JSON(resp)
		}

		userID, err := helper.ExtractUserIDFromClaims(claims)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		companyID, err := helper.ExtractCompanyIDFromClaims(claims)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusUnauthorized).JSON(resp)
		}

		tierLevel, err := helper.ExtractTierLevelFromClaims(claims)
		if err != nil {
			resp := helper.ResponseFailed(err.Error())

			return c.Status(fiber.StatusUnauthorized).JSON(resp)
		}

		c.Locals("userID", userID)
		c.Locals("companyID", companyID)
		c.Locals("tierLevel", tierLevel)

		return c.Next()
	}
}

func AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := os.Getenv("JWT_SECRET_KEY")
		token := c.Cookies("access_token")

		claims, err := helper.ExtractClaimsFromJWT(token, secret)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
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
