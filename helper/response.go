package helper

import "front-office/constant"

type BaseResponseSuccess struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
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
		constant.DataAlreadyExist,
		constant.IncorrectPassword,
		constant.InvalidEmailOrPassword,
		constant.InvalidPassword,
		constant.InvalidPasswordResetLink,
		constant.ConfirmNewPasswordMismatch,
		constant.ConfirmPasswordMismatch:
		statusCode = 400
	case constant.RequestProhibited:
		statusCode = 401
	case constant.DataNotFound, constant.RecordNotFound:
		statusCode = 404
	default:
		statusCode = 500
	}

	resp := ResponseFailed(errorMessage)
	return statusCode, resp
}
