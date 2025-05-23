package helper

import (
	"front-office/common/constant"
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
		constant.ParamSettingIsNotSet:
		statusCode = 400
	case
		constant.WrongCurrentPassword:
		statusCode = 400
		errorMessage = constant.WrongCurrentPassword
	case constant.RequestProhibited,
		constant.TokenExpired,
		constant.UnverifiedUser:
		statusCode = 401
	case constant.DataNotFound,
		constant.RecordNotFound:
		statusCode = 404
		errorMessage = constant.DataNotFound
	case constant.DataAlreadyExist,
		constant.EmailAlreadyExists:
		statusCode = 409
	default:
		statusCode = 500
	}

	resp := ResponseFailed(errorMessage)
	return statusCode, resp
}
