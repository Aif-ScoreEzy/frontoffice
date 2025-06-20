package middleware

import (
	"front-office/internal/apperror"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *apperror.AppError
		if ok := apperror.AsAppError(err, &appErr); ok {
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"message": appErr.Message,
			})
		}

		// Jika error biasa → fallback ke 500
		log.Printf("Unhandled error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}
}
