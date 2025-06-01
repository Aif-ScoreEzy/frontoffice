package helper

import (
	"encoding/json"
	"errors"
	"front-office/common/constant"
	"front-office/common/model"
	"io"
	"net/http"

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
	case constant.UpstreamError:
		statusCode = fiber.StatusBadGateway
	default:
		statusCode = fiber.StatusInternalServerError
	}

	resp := ResponseFailed(errorMessage)
	return statusCode, resp
}

func ParseProCatAPIResponse[T any](response *http.Response) (*model.ProCatAPIResponse[T], error) {
	var apiResponse model.ProCatAPIResponse[T]

	if response == nil {
		return nil, errors.New("nil response")
	}

	dataBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err := json.Unmarshal(dataBytes, &apiResponse); err != nil {
		return nil, err
	}

	apiResponse.StatusCode = response.StatusCode

	return &apiResponse, nil
}
