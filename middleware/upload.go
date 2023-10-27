package middleware

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func FileUpload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		file, err := c.FormFile("image")
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		validExtensions := []string{".jpg", ".jpeg", ".png"}
		ext := filepath.Ext(file.Filename)
		valid := false
		for _, allowedExt := range validExtensions {
			if ext == allowedExt {
				valid = true
				break
			}
		}

		if !valid {
			statusCode, resp := helper.GetError(constant.InvalidImageFile)
			return c.Status(statusCode).JSON(resp)
		}

		if file.Size > 200*1024 {
			statusCode, resp := helper.GetError(constant.FileSizeIsTooLarge)
			return c.Status(statusCode).JSON(resp)
		}

		filename := fmt.Sprintf("%s%s", userID, ext)
		filePath := fmt.Sprintf("./public/%s", filename)

		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				statusCode, resp := helper.GetError(err.Error())
				return c.Status(statusCode).JSON(resp)
			}
		}

		if err := c.SaveFile(file, filePath); err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		c.Locals("filename", filename)

		return c.Next()
	}
}
