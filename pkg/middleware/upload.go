package middleware

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func FileUpload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")

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

		filename := fmt.Sprintf("%d%s", userId, ext)
		filePath := fmt.Sprintf("./storage/uploads/profile/%s", filename)

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

func DocUpload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		fmt.Println("userId: ", userId)

		// get the file upload and type information
		file, err := c.FormFile("file")
		tempType := c.FormValue("tempType")

		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		validExtensions := []string{".csv"}
		ext := filepath.Ext(file.Filename)
		valid := false
		for _, allowedExt := range validExtensions {
			if ext == allowedExt {
				valid = true
				break
			}
		}

		if !valid {
			statusCode, resp := helper.GetError(constant.InvalidDocumentFile)
			return c.Status(statusCode).JSON(resp)
		}

		c.Locals("tempType", tempType)

		return c.Next()
	}
}

func UploadCSVFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		validExtensions := []string{".csv"}
		ext := filepath.Ext(file.Filename)
		valid := false
		for _, allowedExt := range validExtensions {
			if ext == allowedExt {
				valid = true
				break
			}
		}

		if !valid {
			statusCode, resp := helper.GetError(constant.InvalidDocumentFile)
			return c.Status(statusCode).JSON(resp)
		}

		return c.Next()
	}
}
