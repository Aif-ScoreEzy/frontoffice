package middleware

import (
	"front-office/helper"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

func IsRequestValid(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := reflect.New(reflect.TypeOf(model)).Interface()

		if err := c.BodyParser(request); err != nil {
			resp := helper.ResponseFailed("Invalid request format")

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		if errValid := validator.ValidateStruct(request); errValid != nil {
			resp := helper.ResponseFailed(errValid.Error())

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		c.Locals("request", request)

		return c.Next()
	}
}
