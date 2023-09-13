package helper

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

const (
	InvalidPassword   = "password must contain a combination of uppercase, lowercase, number, and symbol"
	IncorrectPassword = "password is incorrect"
	PasswordMismatch  = "please ensure that password and confirm password fields match exactly"
)

func GetError(err error) (int, interface{}) {
	var statusCode int

	switch err.Error() {
	case IncorrectPassword, PasswordMismatch, InvalidPassword:
		statusCode = 400
	default:
		statusCode = 500
	}

	resp := ResponseFailed(err.Error())
	return statusCode, resp
}
