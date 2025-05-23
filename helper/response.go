package helper

import (
	"front-office/common/constant"

	"github.com/gofiber/fiber/v2"
)

type BaseResponseSuccess struct {
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"-"`
}

type BaseResponseFailed struct {
	Message string `json:"message"`
}

func ResponseSuccess(
	message string,
	data interface{},
) BaseResponseSuccess {
	return BaseResponseSuccess{
		Message: message,
		Success: true,
		Data:    data,
	}
}

func ResponseFailed(message string) BaseResponseFailed {
	return BaseResponseFailed{
		Message: message,
	}
}

func GetError(errorMessage string) (int, interface{}) {
	var statusCode int

	switch errorMessage {
	case constant.UserNotFoundForgotEmail:
		statusCode = fiber.StatusOK
	case constant.AlreadyVerified,
		constant.ConfirmNewPasswordMismatch,
		constant.ConfirmPasswordMismatch,
		constant.DuplicateGrading,
		constant.FieldGradingLabelEmpty,
		constant.FieldMinGradeEmpty,
		constant.FieldMaxGradeEmpty,
		constant.FileSizeIsTooLarge,
		constant.IncorrectPassword,
		constant.InvalidActivationLink,
		constant.InvalidStatusValue,
		constant.InvalidDateFormat,
		constant.InvalidEmailOrPassword,
		constant.InvalidImageFile,
		constant.InvalidPassword,
		constant.InvalidPasswordResetLink,
		constant.HeaderTemplateNotValid,
		constant.OnlyUploadCSVfile,
		constant.WrongCurrentPassword,
		constant.ParamSettingIsNotSet:
		statusCode = fiber.StatusBadRequest
	case constant.RequestProhibited,
		constant.TokenExpired,
		constant.UnverifiedUser:
		statusCode = fiber.StatusUnauthorized
	case constant.DataNotFound,
		constant.RecordNotFound:
		statusCode = fiber.StatusNotFound
		errorMessage = constant.DataNotFound
	case constant.TemplateNotFound:
		statusCode = fiber.StatusNotFound
	case constant.DataAlreadyExist,
		constant.EmailAlreadyExists:
		statusCode = fiber.StatusConflict
	default:
		statusCode = fiber.StatusInternalServerError
	}

	resp := ResponseFailed(errorMessage)
	return statusCode, resp
}
